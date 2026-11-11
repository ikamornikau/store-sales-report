# Notes

Index measurements are in [SQL.md](SQL.md). How to run the API and embed the screen is in [README.md](README.md).

## Time

The brief plans for about 3 hours. This took about 5 hours.

## AI tools

Cursor’s agent drafted the seed, the screen, the tests, and these notes. The choices below are the ones to walk through and change live: which rows count, how a store day is bounded, why there are two indexes, and why category pages use `OFFSET`.

## Decisions

Revenue, transaction count, average basket, and items sold use completed transactions only. Refunded and voided rows never add to those figures. The summary still returns `refunded_count` and `voided_count`, and the screen shows them when either is greater than zero. The brief says those rows do not count as revenue; it does not say to hide them.

A day is the store’s calendar day. `from` and `to` are inclusive `YYYY-MM-DD` dates in the store timezone. The bounds are computed on the parameters (`timestamp AT TIME ZONE`), not by wrapping `created_at`, so the index can range-scan. Postgres applies the offset, including DST and `America/St_Johns`. Amounts stay in integer cents in the store currency. There is no conversion.

`from` and `to` may be at most 92 days apart. The server and the client both reject a longer span.

The daily series is zero-filled, so days without sales are 0 and the chart has no gaps.

## Trade-offs

One report is one store and at most 92 days. Two covering indexes (migration `0003`) serve that shape. A partial `WHERE status = 'completed'` index is not used: status is a query parameter, and the generic plan keeps the composite index. Plans and timings are in [SQL.md](SQL.md).

Category pages use `LIMIT`/`OFFSET`. The screen asks for 5 rows. A store has at most six categories, so keyset pagination does not pay for itself. `order_by` and `direction` are an allowlist in Go. `category` is the tie-break, so a page stays stable when the sort key ties.

The API is four `GET` routes under `/api/v1`, snake_case JSON, contract in `backend/docs/openapi.yaml`. Go and TypeScript types are written by hand. An unknown store is 404. A bad date, an inverted range, a span over 92 days, or an invalid sort, limit, or offset is 400 with a plain-text body.

Handlers call the Postgres adapter directly. Validation and the average-basket rounding stay in the handler package. The process uses the standard-library router, JSON logs with a request id, Prometheus metrics, `/health` and `/ready`, read and write timeouts, panic recovery, and graceful shutdown.

The screen is a Vue component. The host passes `api-base-url`, and can pass `fetcher` and `locale` (`en`, `es`, `fr`). Without `locale`, the screen stores the choice in `localStorage`. Copy is typed message objects. The daily chart is chart.js. The three report calls run in parallel and the in-flight batch is aborted when the filter changes. `pnpm build:lib` publishes the screen with `vue`, `chart.js`, and `vue-chartjs` as peers.

## Schema

The given columns stay. Changes:

- `CHECK` lists for `currency`, `category`, `payment_method`, and `status`, so a new value is a migration.
- `items` is `items_count`, with `items_count > 0` and `total_cents >= 0`. The brief’s `items` is a count.
- `BIGINT GENERATED ALWAYS AS IDENTITY` instead of `BIGSERIAL`.
- A unique index on `lower(name)`.
- `is_timezone()` checks the timezone against `pg_timezone_names`, so a typo cannot move every day boundary.

There is no `local_date` column. The local day is `AT TIME ZONE` at read time, so stored rows do not depend on a store’s timezone staying the same. The report indexes are a later migration, so the before and after plans in [SQL.md](SQL.md) are the same queries.

## Seed

Compose loads 50 stores and exactly 2,000,000 transactions over the 90 days ending at `now()`. The status draw is about 2% refunded and 1% voided. Nine stores are CAD, in `America/St_Johns`, `America/Toronto`, and `America/Vancouver`. Volume is skewed: the busiest store is about 4.3× the quietest. `setseed` makes the draw repeatable, and the script raises if a check fails. The seed image does not reload when `stores` or `transactions` already has rows. On a local machine the load finished in about 20 seconds.

## Next

- Authentication between services. The local API is open so Compose can serve the demo.
- Partition `transactions` by month if the table grows well past this seed. The current indexes are sized for one store and 92 days.
