-- One store and a created_at window. Summary reads every status, and
-- COUNT(id) needs id in the index for an index-only scan.
CREATE INDEX transactions_store_created_at_idx
    ON transactions (store_id, created_at)
    INCLUDE (status, total_cents, items_count, id);

-- Daily and category revenue keep completed sales only.
-- category is included so that rollup does not visit the heap.
CREATE INDEX transactions_store_status_created_at_idx
    ON transactions (store_id, status, created_at)
    INCLUDE (total_cents, category);
