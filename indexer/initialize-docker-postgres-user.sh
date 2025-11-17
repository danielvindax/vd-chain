#!/bin/bash

# Initialize the postgres docker image.

set -e
set -u

# Initialize extensions locally that are useful for development.
initialize_extensions() {
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
create extension if not exists "plpgsql";
create extension if not exists "plpgsql_check";
create extension if not exists "plprofiler";
create extension if not exists "pldbgapi";
EOSQL
}

initialize_extensions
