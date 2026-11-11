-- 50 stores and 2_000_000 transactions over the 90 days ending at now().
-- About 3% of transactions are refunded (2%) or voided (1%).
-- CAD: America/St_Johns, America/Toronto, America/Vancouver.
-- USD: America/New_York, America/Chicago, America/Denver,
--      America/Los_Angeles, America/Anchorage, Pacific/Honolulu.
--
-- psql "postgres://store:store@localhost:5432/store_sales_report?sslmode=disable" -v ON_ERROR_STOP=1 -f seed.sql

SELECT (
    EXISTS (SELECT 1 FROM stores)
    OR EXISTS (SELECT 1 FROM transactions)
)::integer AS skip_seed
\gset

\if :skip_seed
\echo seed skipped: tables already contain data
\quit
\endif

BEGIN;

SET LOCAL synchronous_commit TO off;
SET LOCAL jit TO off;
SET LOCAL max_parallel_workers_per_gather TO 0;
SET LOCAL timezone TO 'UTC';

TRUNCATE TABLE transactions, stores RESTART IDENTITY CASCADE;

INSERT INTO stores (id, name, currency, timezone) VALUES
    -- Canada: Newfoundland, Eastern, Pacific
    ('S-1001', 'St. John''s Harbour Market', 'CAD', 'America/St_Johns'),
    ('S-1002', 'Corner Brook Convenience', 'CAD', 'America/St_Johns'),
    ('S-1003', 'Gander Crossroads Market', 'CAD', 'America/St_Johns'),
    ('S-1004', 'Toronto Queen Street', 'CAD', 'America/Toronto'),
    ('S-1005', 'Ottawa Rideau Market', 'CAD', 'America/Toronto'),
    ('S-1006', 'Montreal Plateau Shop', 'CAD', 'America/Toronto'),
    ('S-1007', 'Vancouver Robson Market', 'CAD', 'America/Vancouver'),
    ('S-1008', 'Victoria Inner Harbour', 'CAD', 'America/Vancouver'),
    ('S-1009', 'Surrey Central Store', 'CAD', 'America/Vancouver'),

    -- United States: Eastern
    ('S-1010', 'Manhattan Deli', 'USD', 'America/New_York'),
    ('S-1011', 'Brooklyn Corner Market', 'USD', 'America/New_York'),
    ('S-1012', 'Boston Harbor Shop', 'USD', 'America/New_York'),
    ('S-1013', 'Miami Beach Market', 'USD', 'America/New_York'),
    ('S-1014', 'Atlanta Peachtree Store', 'USD', 'America/New_York'),
    ('S-1015', 'Philadelphia Market Street', 'USD', 'America/New_York'),
    ('S-1016', 'Charlotte Uptown Market', 'USD', 'America/New_York'),
    ('S-1017', 'Washington Capitol Shop', 'USD', 'America/New_York'),
    ('S-1018', 'Jacksonville Beach Store', 'USD', 'America/New_York'),
    ('S-1019', 'Tampa Bay Market', 'USD', 'America/New_York'),
    ('S-1020', 'Orlando Lake Market', 'USD', 'America/New_York'),
    ('S-1021', 'Newark Ironbound Shop', 'USD', 'America/New_York'),

    -- United States: Central
    ('S-1022', 'Chicago Loop Market', 'USD', 'America/Chicago'),
    ('S-1023', 'Dallas Uptown Store', 'USD', 'America/Chicago'),
    ('S-1024', 'Houston Heights Market', 'USD', 'America/Chicago'),
    ('S-1025', 'Minneapolis Nicollet Shop', 'USD', 'America/Chicago'),
    ('S-1026', 'Austin Congress Market', 'USD', 'America/Chicago'),
    ('S-1027', 'Nashville Broadway Store', 'USD', 'America/Chicago'),
    ('S-1028', 'Milwaukee Lakefront Market', 'USD', 'America/Chicago'),
    ('S-1029', 'Kansas City Westport', 'USD', 'America/Chicago'),
    ('S-1030', 'New Orleans Magazine Shop', 'USD', 'America/Chicago'),
    ('S-1031', 'San Antonio River Market', 'USD', 'America/Chicago'),

    -- United States: Mountain
    ('S-1032', 'Denver LoDo Market', 'USD', 'America/Denver'),
    ('S-1033', 'Salt Lake City Temple Shop', 'USD', 'America/Denver'),
    ('S-1034', 'Albuquerque Old Town', 'USD', 'America/Denver'),
    ('S-1035', 'Boulder Pearl Street', 'USD', 'America/Denver'),
    ('S-1036', 'Billings Rimrock Store', 'USD', 'America/Denver'),

    -- United States: Pacific
    ('S-1037', 'Los Angeles Sunset Market', 'USD', 'America/Los_Angeles'),
    ('S-1038', 'Seattle Pike Place Shop', 'USD', 'America/Los_Angeles'),
    ('S-1039', 'San Francisco Mission Store', 'USD', 'America/Los_Angeles'),
    ('S-1040', 'Portland Hawthorne Market', 'USD', 'America/Los_Angeles'),
    ('S-1041', 'San Diego Gaslamp Shop', 'USD', 'America/Los_Angeles'),
    ('S-1042', 'Las Vegas Strip Market', 'USD', 'America/Los_Angeles'),
    ('S-1043', 'Sacramento Midtown Store', 'USD', 'America/Los_Angeles'),
    ('S-1044', 'San Jose Willow Glen', 'USD', 'America/Los_Angeles'),

    -- United States: Alaska
    ('S-1045', 'Anchorage Downtown Market', 'USD', 'America/Anchorage'),
    ('S-1046', 'Fairbanks Geist Store', 'USD', 'America/Anchorage'),
    ('S-1047', 'Juneau Waterfront Shop', 'USD', 'America/Anchorage'),

    -- United States: Hawaii
    ('S-1048', 'Honolulu Waikiki Market', 'USD', 'Pacific/Honolulu'),
    ('S-1049', 'Hilo Bayfront Store', 'USD', 'Pacific/Honolulu'),
    ('S-1050', 'Kahului Harbor Shop', 'USD', 'Pacific/Honolulu');

SELECT setseed(0.20261111);

CREATE TEMP TABLE seed_store_slots (
    n integer PRIMARY KEY,
    id text NOT NULL
) ON COMMIT DROP;

INSERT INTO seed_store_slots (n, id)
SELECT (row_number() OVER (ORDER BY id) - 1)::integer, id
FROM stores;

INSERT INTO transactions (
    store_id,
    created_at,
    category,
    payment_method,
    items_count,
    total_cents,
    status
)
SELECT
    s.id,
    now() - (u.r_day * interval '90 days'),
    a.category,
    a.payment_method,
    a.items_count,
    (
        a.items_count * CASE a.category
            WHEN 'beverages' THEN 150 + floor(u.r_price * 550)::integer
            WHEN 'snacks' THEN 99 + floor(u.r_price * 400)::integer
            WHEN 'tobacco' THEN 800 + floor(u.r_price * 1000)::integer
            WHEN 'grocery' THEN 200 + floor(u.r_price * 1300)::integer
            WHEN 'lottery' THEN 200 + floor(u.r_price * 800)::integer
            WHEN 'household' THEN 250 + floor(u.r_price * 1750)::integer
        END
    )::bigint,
    a.status
FROM generate_series(1, 2000000) AS g(i)
JOIN seed_store_slots AS s ON s.n = g.i % 50
-- Tie the draw to g.i. An uncorrelated subquery is planned once, so random()
-- would otherwise return the same value for every transaction.
CROSS JOIN LATERAL (
    SELECT
        random() AS r_day,
        random() AS r_cat,
        random() AS r_pay,
        random() AS r_status,
        random() AS r_price,
        random() AS r_items
    FROM (SELECT g.i) AS tie(i)
) AS u
CROSS JOIN LATERAL (
    SELECT
        CASE
            WHEN u.r_cat < 0.28 THEN 'beverages'
            WHEN u.r_cat < 0.52 THEN 'snacks'
            WHEN u.r_cat < 0.66 THEN 'tobacco'
            WHEN u.r_cat < 0.84 THEN 'grocery'
            WHEN u.r_cat < 0.92 THEN 'lottery'
            ELSE 'household'
        END AS category,
        CASE
            WHEN u.r_pay < 0.52 THEN 'card'
            WHEN u.r_pay < 0.80 THEN 'cash'
            WHEN u.r_pay < 0.92 THEN 'mobile'
            ELSE 'ebt'
        END AS payment_method,
        CASE
            WHEN u.r_status < 0.02 THEN 'refunded'
            WHEN u.r_status < 0.03 THEN 'voided'
            ELSE 'completed'
        END AS status,
        1 + LEAST(floor(ln(greatest(u.r_items, 1e-9)) / ln(0.55))::integer, 11) AS items_count
) AS a;

DO $$
DECLARE
    store_count integer;
    cad_zones integer;
    usd_zones integer;
    tx_count bigint;
    refunded_ratio numeric;
    voided_ratio numeric;
    oldest timestamptz;
    newest timestamptz;
BEGIN
    SELECT count(*) INTO store_count FROM stores;

    SELECT count(DISTINCT timezone) INTO cad_zones
    FROM stores
    WHERE currency = 'CAD';

    SELECT count(DISTINCT timezone) INTO usd_zones
    FROM stores
    WHERE currency = 'USD';

    SELECT
        count(*),
        count(*) FILTER (WHERE status = 'refunded')::numeric / count(*),
        count(*) FILTER (WHERE status = 'voided')::numeric / count(*),
        min(created_at),
        max(created_at)
    INTO tx_count, refunded_ratio, voided_ratio, oldest, newest
    FROM transactions;

    IF store_count <> 50
        OR cad_zones <> 3
        OR usd_zones <> 6
        OR EXISTS (
            SELECT 1
            FROM stores
            WHERE currency = 'CAD'
              AND timezone NOT IN ('America/St_Johns', 'America/Toronto', 'America/Vancouver')
        )
        OR EXISTS (
            SELECT 1
            FROM stores
            WHERE currency = 'USD'
              AND timezone NOT IN (
                  'America/New_York',
                  'America/Chicago',
                  'America/Denver',
                  'America/Los_Angeles',
                  'America/Anchorage',
                  'Pacific/Honolulu'
              )
        )
        OR tx_count <> 2000000
        OR refunded_ratio + voided_ratio < 0.025
        OR refunded_ratio + voided_ratio > 0.035
        OR refunded_ratio < 0.015
        OR voided_ratio < 0.005
        OR newest > now()
        OR oldest < now() - interval '90 days' - interval '1 second'
    THEN
        RAISE EXCEPTION
            'seed check failed: stores=% cad_zones=% usd_zones=% tx=% refunded=% voided=% oldest=% newest=%',
            store_count,
            cad_zones,
            usd_zones,
            tx_count,
            refunded_ratio,
            voided_ratio,
            oldest,
            newest;
    END IF;

    RAISE NOTICE
        'seeded % stores and % transactions; refunded %, voided %',
        store_count,
        tx_count,
        round(refunded_ratio, 4),
        round(voided_ratio, 4);
END $$;

ANALYZE stores;
ANALYZE transactions;

COMMIT;
