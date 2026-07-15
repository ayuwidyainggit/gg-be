#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create Citus extension
    CREATE EXTENSION IF NOT EXISTS citus;
    
    -- Create Citus Columnar extension
    CREATE EXTENSION IF NOT EXISTS citus_columnar;
    
    -- Show installed extensions
    \dx
EOSQL

echo "Citus extensions initialized successfully!"

