-- 50 stores and 2_000_000 transactions over the 90 days ending at now().
-- Store volume follows a smooth curve, about 4.3× from the quietest market
-- to the busiest. The same margin and a fixed cost at the store ranked
-- 37th of 50 put operating profit near 37% / 62% / -7% for the top 10%,
-- the top 20%, and the bottom quartile.
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

-- Rank 1 is the busiest store. Big-city and tourist locations sit at the
-- top; small-town markets sit at the bottom. Weight for rank r is
-- 5.2 * (1.22 / 5.2) ^ (((r - 1) / 49) ^ 0.7). Largest remainder turns
-- those weights into quotas that sum to exactly 2_000_000.
CREATE TEMP TABLE seed_store_slots (
    n integer PRIMARY KEY,
    id text NOT NULL,
    tx_count integer NOT NULL
) ON COMMIT DROP;

INSERT INTO seed_store_slots (n, id, tx_count)
WITH traffic AS (
    SELECT *
    FROM (VALUES
        (1, 'S-1010'),  -- Manhattan Deli
        (2, 'S-1022'),  -- Chicago Loop Market
        (3, 'S-1037'),  -- Los Angeles Sunset Market
        (4, 'S-1042'),  -- Las Vegas Strip Market
        (5, 'S-1004'),  -- Toronto Queen Street
        (6, 'S-1048'),  -- Honolulu Waikiki Market
        (7, 'S-1007'),  -- Vancouver Robson Market
        (8, 'S-1013'),  -- Miami Beach Market
        (9, 'S-1039'),  -- San Francisco Mission Store
        (10, 'S-1012'), -- Boston Harbor Shop
        (11, 'S-1011'), -- Brooklyn Corner Market
        (12, 'S-1038'), -- Seattle Pike Place Shop
        (13, 'S-1017'), -- Washington Capitol Shop
        (14, 'S-1015'), -- Philadelphia Market Street
        (15, 'S-1020'), -- Orlando Lake Market
        (16, 'S-1006'), -- Montreal Plateau Shop
        (17, 'S-1014'), -- Atlanta Peachtree Store
        (18, 'S-1023'), -- Dallas Uptown Store
        (19, 'S-1024'), -- Houston Heights Market
        (20, 'S-1027'), -- Nashville Broadway Store
        (21, 'S-1026'), -- Austin Congress Market
        (22, 'S-1032'), -- Denver LoDo Market
        (23, 'S-1041'), -- San Diego Gaslamp Shop
        (24, 'S-1030'), -- New Orleans Magazine Shop
        (25, 'S-1005'), -- Ottawa Rideau Market
        (26, 'S-1044'), -- San Jose Willow Glen
        (27, 'S-1040'), -- Portland Hawthorne Market
        (28, 'S-1019'), -- Tampa Bay Market
        (29, 'S-1025'), -- Minneapolis Nicollet Shop
        (30, 'S-1016'), -- Charlotte Uptown Market
        (31, 'S-1031'), -- San Antonio River Market
        (32, 'S-1008'), -- Victoria Inner Harbour
        (33, 'S-1043'), -- Sacramento Midtown Store
        (34, 'S-1028'), -- Milwaukee Lakefront Market
        (35, 'S-1009'), -- Surrey Central Store
        (36, 'S-1021'), -- Newark Ironbound Shop
        (37, 'S-1033'), -- Salt Lake City Temple Shop
        (38, 'S-1029'), -- Kansas City Westport
        (39, 'S-1018'), -- Jacksonville Beach Store
        (40, 'S-1035'), -- Boulder Pearl Street
        (41, 'S-1034'), -- Albuquerque Old Town
        (42, 'S-1045'), -- Anchorage Downtown Market
        (43, 'S-1001'), -- St. John's Harbour Market
        (44, 'S-1050'), -- Kahului Harbor Shop
        (45, 'S-1047'), -- Juneau Waterfront Shop
        (46, 'S-1049'), -- Hilo Bayfront Store
        (47, 'S-1046'), -- Fairbanks Geist Store
        (48, 'S-1036'), -- Billings Rimrock Store
        (49, 'S-1002'), -- Corner Brook Convenience
        (50, 'S-1003')  -- Gander Crossroads Market
    ) AS v(rank, id)
),
weighted AS (
    SELECT
        rank,
        id,
        5.2 * power(1.22 / 5.2, power((rank - 1)::numeric / 49, 0.7)) AS weight
    FROM traffic
),
quotas AS (
    SELECT
        rank,
        id,
        2000000 * weight / sum(weight) OVER () AS exact_count
    FROM weighted
),
apportioned AS (
    SELECT
        rank,
        id,
        floor(exact_count)::integer AS base_count,
        exact_count - floor(exact_count) AS remainder,
        (2000000 - sum(floor(exact_count)) OVER ())::integer AS leftover
    FROM quotas
)
SELECT
    rank,
    id,
    base_count + CASE
        WHEN row_number() OVER (ORDER BY remainder DESC, id) <= leftover THEN 1
        ELSE 0
    END
FROM apportioned;

DO $$
BEGIN
    IF (SELECT count(*) FROM seed_store_slots) <> 50
        OR (SELECT count(*) FROM seed_store_slots AS slot JOIN stores AS store ON store.id = slot.id) <> 50
        OR (SELECT sum(tx_count) FROM seed_store_slots) <> 2000000
        OR (SELECT min(tx_count) FROM seed_store_slots) <= 0
    THEN
        RAISE EXCEPTION 'store traffic quotas do not cover 50 stores and 2000000 transactions';
    END IF;
END $$;

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
FROM seed_store_slots AS s
CROSS JOIN LATERAL generate_series(1, s.tx_count) AS g(i)
-- Tie the draw to the store and g.i. An uncorrelated subquery is planned
-- once, so random() would otherwise return the same value for every
-- transaction. g.i repeats across stores, so the store slot has to be part
-- of the tie.
CROSS JOIN LATERAL (
    SELECT
        random() AS r_day,
        random() AS r_cat,
        random() AS r_pay,
        random() AS r_status,
        random() AS r_price,
        random() AS r_items
    FROM (SELECT s.n, g.i) AS tie(n, i)
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
    traffic_ratio numeric;
    top_decile_profit numeric;
    top_quintile_profit numeric;
    bottom_quartile_profit numeric;
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

    -- Ticket mix does not depend on the store, so transaction counts stand
    -- in for revenue. Fixed cost is the volume of the store ranked 37th,
    -- just above the bottom quartile. Profit is volume minus that cost.
    WITH volume AS (
        SELECT count(*)::numeric AS tx
        FROM transactions
        GROUP BY store_id
    ),
    ranked AS (
        SELECT tx, row_number() OVER (ORDER BY tx DESC) AS rk
        FROM volume
    ),
    fixed AS (
        SELECT tx AS cost
        FROM ranked
        WHERE rk = 37
    ),
    profit AS (
        SELECT ranked.rk, ranked.tx, ranked.tx - fixed.cost AS profit
        FROM ranked
        CROSS JOIN fixed
    )
    SELECT
        max(tx) / min(tx),
        sum(profit) FILTER (WHERE rk <= 5) / nullif(sum(profit), 0),
        sum(profit) FILTER (WHERE rk <= 10) / nullif(sum(profit), 0),
        sum(profit) FILTER (WHERE rk > 37) / nullif(sum(profit), 0)
    INTO traffic_ratio, top_decile_profit, top_quintile_profit, bottom_quartile_profit
    FROM profit;

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
        OR traffic_ratio IS NULL
        OR traffic_ratio < 4.0
        OR traffic_ratio > 4.6
        OR top_decile_profit IS NULL
        OR top_decile_profit < 0.33
        OR top_decile_profit > 0.42
        OR top_quintile_profit IS NULL
        OR top_quintile_profit < 0.56
        OR top_quintile_profit > 0.68
        OR bottom_quartile_profit IS NULL
        OR bottom_quartile_profit < -0.12
        OR bottom_quartile_profit > -0.03
    THEN
        RAISE EXCEPTION
            'seed check failed: stores=% cad_zones=% usd_zones=% tx=% refunded=% voided=% oldest=% newest=% traffic_ratio=% top10_profit=% top20_profit=% bottom_quartile_profit=%',
            store_count,
            cad_zones,
            usd_zones,
            tx_count,
            refunded_ratio,
            voided_ratio,
            oldest,
            newest,
            traffic_ratio,
            top_decile_profit,
            top_quintile_profit,
            bottom_quartile_profit;
    END IF;

    RAISE NOTICE
        'seeded % stores and % transactions; refunded %, voided %; profit top 10%% %, top 20%% %, bottom quartile %',
        store_count,
        tx_count,
        round(refunded_ratio, 4),
        round(voided_ratio, 4),
        round(top_decile_profit, 3),
        round(top_quintile_profit, 3),
        round(bottom_quartile_profit, 3);
END $$;

ANALYZE stores;
ANALYZE transactions;

COMMIT;
