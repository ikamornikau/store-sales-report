CREATE FUNCTION is_timezone(tz TEXT) RETURNS BOOLEAN
LANGUAGE SQL
STABLE
PARALLEL SAFE
SET search_path = pg_catalog
AS $$
    SELECT EXISTS (
        SELECT 1
        FROM pg_timezone_names
        WHERE name = tz
    );
$$;

CREATE TABLE stores (
    id TEXT PRIMARY KEY, -- e.g. 'S-1042'
    name TEXT NOT NULL,
    currency CHAR(3) NOT NULL CHECK (currency IN ('USD', 'CAD')),
    timezone TEXT NOT NULL CHECK (is_timezone(timezone)) -- e.g. 'America/New_York'
);

CREATE UNIQUE INDEX stores_name_lower_unique_idx ON stores (lower(name));
