#!/bin/sh

until pg_isready -h "${NX_POSTGRES_USER_DB_HOST}" -p "${NX_POSTGRES_USER_DB_PORT}" -U "${NX_POSTGRES_USER_DB_USER}" -d "${NX_POSTGRES_USER_DB_NAME}";
do
    echo >&2 "Postgres is unavailable - sleeping"
    sleep 5
done

until nc -z "${NX_CLICKHOUSE_HOST}" "${NX_CLICKHOUSE_PORT}";
do
    echo >&2 "Clickhouse is unavailable - sleeping"
    sleep 5
done

postgresDatabaseURL="postgres://${NX_POSTGRES_USER_DB_USER}:${NX_POSTGRES_USER_DB_PASSWORD}@${NX_POSTGRES_USER_DB_HOST}:${NX_POSTGRES_USER_DB_PORT}/${NX_POSTGRES_USER_DB_NAME}?sslmode=${NX_POSTGRES_USER_DB_SSLMODE}"
echo "run database migrations"

echo "postgres database url: $postgresDatabaseURL"
/usr/local/bin/migrate -verbose -source file:///usr/local/postgresql -database "$postgresDatabaseURL" up
if [ ! $? -eq 0 ]; then
    echo "postgres database migration failed, exiting"
    exit 1
fi

clickhouseDatabaseURL="clickhouse://${NX_CLICKHOUSE_HOST}:${NX_CLICKHOUSE_PORT}?database=${NX_CLICKHOUSE_DB}&username=${NX_CLICKHOUSE_USER}&password=${NX_CLICKHOUSE_PASSWORD}&x-multi-statement=true"

echo "clickhouse database url: $clickhouseDatabaseURL"
# /usr/local/bin/migrate -verbose -source file:///usr/local/clickhouse -database "$clickhouseDatabaseURL" up
# if [ ! $? -eq 0 ]; then
#     echo "clickhouse database migration failed, exiting"
#     exit 1
# fi

echo "run database migrations, complete"

echo "starting numexa auth service..."
exec "$@"
