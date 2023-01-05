#!/bin/sh

set -e

echo "run db migration"

# Other option is do add the wait-script here
# ./wait-for.sh "db:5432"
# source /app/app.env
# /app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start app"
# Run all arguments passed to the script
exec "$@"