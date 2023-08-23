#!/bin/sh

until nc -z "${NX_CLICKHOUSE_HOST}" "${NX_CLICKHOUSE_PORT}";
do
    echo >&2 "Clickhouse is unavailable - sleeping"
    sleep 5
done


echo "starting numexa vibe service..."
exec "$@"