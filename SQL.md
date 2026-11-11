# Report indexes

Migration `0003_transactions_report_indexes` adds two indexes on `transactions`. Before that migration the table had only its primary key. Measurements below are from the seeded database: 50 stores, 2,000,000 transactions, heap 188 MB.

Server: PostgreSQL 18.6, `shared_buffers` 128MB, `work_mem` 4MB, `max_parallel_workers_per_gather` 2. `track_io_timing` was on for these plans.

The three report queries are `GetTransactionSummary`, `GetDailyRevenue`, and `ListCategoryRevenue` in `backend/internal/adapters/postgres/reports.go`. Each one keeps a single store and a `created_at` window. The window is the store's local calendar days, so the bounds are `($date::timestamp AT TIME ZONE store.timezone)` and cannot be constants.

The slow case is the busiest store over the whole seeded history. `S-1010` (America/New_York) has 86,332 transactions, from 2026-08-13 through 2026-11-11. That span is inside the 92-day limit the UI allows.

## Indexes

`transactions_store_created_at_idx` on `(store_id, created_at) INCLUDE (status, total_cents, items_count, id)`.

The summary reads every status in the window: completed revenue, plus refunded and voided counts. `status` is not an equality predicate, so it stays in `INCLUDE`. The aggregate uses `total_cents`, `items_count`, and `COUNT(id)`. With `id` in the index the plan is an index-only scan (`Heap Fetches: 0`) instead of a bitmap heap scan of those 86,332 rows.

`transactions_store_status_created_at_idx` on `(store_id, status, created_at) INCLUDE (total_cents, category)`.

Daily and category revenue keep `status = completed` only. `status` is the second key, so refunded and voided rows (about 3% of the table) are not visited. `total_cents` covers both sums. `category` covers the category `GROUP BY`, so that query is index-only too.

A partial index `WHERE status = 'completed'` would miss the generic plan. The API passes status as `$4`, and `EXPLAIN (GENERIC_PLAN)` on that SQL chooses this composite index.

## Slowest query

`GetTransactionSummary` for `S-1010`, 2026-08-13 through 2026-11-11. On the unindexed table it was the slowest of the three:

| Query | Before | After |
| --- | ---: | ---: |
| Summary | 88.143 ms | 19.246 ms |
| Daily revenue | 67.810 ms | 26.531 ms |
| Category revenue | 53.377 ms | 14.187 ms |

The after times are the second execution, with the index pages already in `shared_buffers`. The first summary execution after the migration was 33.512 ms and read 810 blocks (7.954 ms of I/O).

Without an index the summary parallel-scans all 2,000,000 rows and throws away every store except `S-1010` (637,889 rows removed in each of 3 workers). The `created_at` predicate is only a join filter, applied after that scan. With the index it reads the 86,332 index entries for that store and does not touch the heap.

### Before

```text
Aggregate  (cost=46190.18..46190.19 rows=1 width=40) (actual time=86.709..88.047 rows=1.00 loops=1)
  Buffers: shared hit=15866 read=8258
  I/O Timings: shared read=7.659
  ->  Nested Loop Left Join  (cost=1000.00..45949.25 rows=9637 width=29) (actual time=32.273..78.976 rows=86332.00 loops=1)
        Join Filter: ((t.created_at >= ('2026-08-13 00:00:00'::timestamp without time zone AT TIME ZONE s.timezone)) AND (t.created_at < ('2026-11-12 00:00:00'::timestamp without time zone AT TIME ZONE s.timezone)))
        Buffers: shared hit=15866 read=8258
        I/O Timings: shared read=7.659
        ->  Seq Scan on stores s  (cost=0.00..1.62 rows=1 width=24) (actual time=0.016..0.019 rows=1.00 loops=1)
              Filter: (id = 'S-1010'::text)
              Rows Removed by Filter: 49
              Buffers: shared hit=1
        ->  Gather  (cost=1000.00..44212.97 rows=86733 width=44) (actual time=31.984..41.189 rows=86332.00 loops=1)
              Workers Planned: 2
              Workers Launched: 2
              Buffers: shared hit=15865 read=8258
              I/O Timings: shared read=7.659
              ->  Parallel Seq Scan on transactions t  (cost=0.00..34539.67 rows=36139 width=44) (actual time=28.662..35.865 rows=28777.33 loops=3)
                    Filter: (store_id = 'S-1010'::text)
                    Rows Removed by Filter: 637889
                    Buffers: shared hit=15865 read=8258
                    I/O Timings: shared read=7.659
Planning Time: 2.230 ms
Execution Time: 88.143 ms
```

### After

```text
Aggregate  (cost=936.43..936.44 rows=1 width=40) (actual time=19.187..19.188 rows=1.00 loops=1)
  Buffers: shared hit=812
  ->  Nested Loop Left Join  (cost=0.43..693.65 rows=9711 width=29) (actual time=0.023..10.614 rows=86332.00 loops=1)
        Buffers: shared hit=812
        ->  Seq Scan on stores s  (cost=0.00..1.62 rows=1 width=24) (actual time=0.005..0.009 rows=1.00 loops=1)
              Filter: (id = 'S-1010'::text)
              Rows Removed by Filter: 49
              Buffers: shared hit=1
        ->  Index Only Scan using transactions_store_created_at_idx on transactions t  (cost=0.43..594.92 rows=9711 width=44) (actual time=0.013..5.624 rows=86332.00 loops=1)
              Index Cond: ((store_id = 'S-1010'::text) AND (created_at >= ('2026-08-13 00:00:00'::timestamp without time zone AT TIME ZONE s.timezone)) AND (created_at < ('2026-11-12 00:00:00'::timestamp without time zone AT TIME ZONE s.timezone)))
              Heap Fetches: 0
              Index Searches: 1
              Buffers: shared hit=811
Planning Time: 0.103 ms
Execution Time: 19.246 ms
```

Daily and category revenue use `transactions_store_status_created_at_idx` the same way: index-only scan, `Heap Fetches: 0`, index condition `(store_id, status, created_at)`.
