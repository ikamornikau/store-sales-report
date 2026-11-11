CREATE TABLE transactions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    store_id TEXT NOT NULL REFERENCES stores(id),
    created_at TIMESTAMPTZ NOT NULL,
    category TEXT NOT NULL CHECK (category IN (
        'beverages', 'snacks', 'tobacco', 'grocery', 'lottery', 'household')
    ),
    payment_method TEXT NOT NULL CHECK (payment_method IN ('cash', 'card', 'ebt', 'mobile')),
    items_count INT NOT NULL CHECK (items_count > 0),
    total_cents BIGINT NOT NULL CHECK (total_cents >= 0),
    status TEXT NOT NULL CHECK (status IN ('completed', 'refunded', 'voided'))
);
