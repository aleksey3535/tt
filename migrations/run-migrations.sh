#!/bin/sh

echo "⏳ Waiting for DB to be ready..."

for i in $(seq 1 10); do
  migrate -path=/migrations -database "$MIGRATE_DATABASE_URL" up && exit 0
  echo "Retrying in 3s..."
  sleep 3
done

echo "❌ Migration failed after retries"
exit 1
