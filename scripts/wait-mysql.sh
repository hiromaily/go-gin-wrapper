#!/bin/sh
# wait_mysql.sh

#command: ["./wait_mysql.sh", "mysql", "/go/bin/api", "--host", "0.0.0.0", "--port", "8081", "--config", "config/example.yml"]
set -e

host="$1"
shift
cmd="$@"

while ! mysqladmin ping -h"$host" --silent; do
    >&2 echo "Database is unavailable - sleeping"
    sleep 1
done

>&2 echo "Database is up - executing command"
exec $cmd

