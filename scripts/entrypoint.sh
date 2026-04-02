#!/bin/sh
set -e

echo "[entrypoint] Aguardando MySQL em ${DB_HOST}:${DB_PORT}..."
until mysqladmin ping -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" --silent --default-auth=mysql_native_password 2>/dev/null; do
    echo "[entrypoint] MySQL nao disponivel ainda - aguardando 2s..."
    sleep 2
done
echo "[entrypoint] MySQL disponivel!"

echo "[entrypoint] Aplicando migrations..."
mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" --default-auth=mysql_native_password < /app/migrations/001_create_orders.sql
echo "[entrypoint] Migrations aplicadas!"

echo "[entrypoint] Iniciando aplicacao..."
exec /app/server
