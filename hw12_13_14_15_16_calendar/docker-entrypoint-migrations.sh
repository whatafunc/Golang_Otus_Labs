#!/bin/sh

# Exit immediately if a command exits with a non-zero status.
set -e

# Check if POSTGRES_DSN is set and not empty
if [ -n "$POSTGRES_DSN" ]; then
  echo "POSTGRES_DSN is set. Running database migrations..."

  # Parse the DSN string to build a URL for goose
  DB_USER=$(echo "$POSTGRES_DSN" | sed -n 's/.*user=\([^ ]*\).*/\1/p')
  DB_PASSWORD=$(echo "$POSTGRES_DSN" | sed -n 's/.*password=\([^ ]*\).*/\1/p')
  DB_HOST=$(echo "$POSTGRES_DSN" | sed -n 's/.*host=\([^ ]*\).*/\1/p')
  DB_PORT=$(echo "$POSTGRES_DSN" | sed -n 's/.*port=\([^ ]*\).*/\1/p')
  DB_NAME=$(echo "$POSTGRES_DSN" | sed -n 's/.*dbname=\([^ ]*\).*/\1/p')
  SSL_MODE=$(echo "$POSTGRES_DSN" | sed -n 's/.*sslmode=\([^ ]*\).*/\1/p')

  # Construct the URL for goose
  DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${SSL_MODE}"

  # Run migrations.
  goose -dir /migrations postgres "$DATABASE_URL" up

  echo "Migrations complete."
else
  echo "POSTGRES_DSN not set, skipping migrations."
fi

echo "Starting application..."
# Execute the command passed as arguments to this script (the Dockerfile's CMD)
exec "$@"
