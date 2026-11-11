FROM postgres:18-alpine

COPY seed.sql /seed.sql

ENTRYPOINT ["psql", "-v", "ON_ERROR_STOP=1", "-f", "/seed.sql"]
