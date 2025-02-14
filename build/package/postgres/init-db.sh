#!/bin/sh
set -e

POSTGRES_HOST=${POSTGRES_HOST:-postgres}

echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h "$POSTGRES_HOST" -p 5432 -U "$POSTGRES_USER"; do
  sleep 2
done

echo "PostgreSQL is ready. Checking if database $POSTGRES_DB exists..."
DB_EXISTS=$(psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -tAc "SELECT 1 FROM pg_database WHERE datname='$POSTGRES_DB';")

if [ "$DB_EXISTS" != "1" ]; then
  echo "Database $POSTGRES_DB does not exist. Creating..."
  psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -c "CREATE DATABASE $POSTGRES_DB;"
  psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
  echo "Database $POSTGRES_DB and extension uuid-ossp created successfully."
else
  echo "Database $POSTGRES_DB already exists. Skipping creation."
fi

echo "PostgreSQL setup complete."
exit 0

