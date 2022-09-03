#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    ALTER DATABASE ${POSTGRES_DB} SET search_path TO public;
EOSQL

psql -v ON_ERROR_STOP=1 -v --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f ./migrations/000001_init.up.sql
