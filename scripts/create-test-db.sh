#!/bin/sh
###############################################################################
# create test database
###############################################################################

# settings
DB_HOST=127.0.0.1
#DB_USER=root
#DB_PASS=
#DB_PORT=3306

echo ${DB_USER}
echo ${DB_PASS}
echo ${DB_PORT}
echo ${DB_NAME}

SQL_DIR=${GOPATH}/src/github.com/hiromaily/go-gin-wrapper/test/sql

# Dump
# mysqldump -u root -p hiromaily > data_hiromaily.sql


# create test db
# Note:
# mysql command is required to run
# `brew services list`
# `brew services start mysql`
# if facing `sh: mysql: command not found`, try to run `brew link mysql` or `brew link --overwrite mysql`
expect -c "
    set timeout 30
    spawn sh -c \"mysql -u${DB_USER:-root} -p -h${DB_HOST} -P${DB_PORT:-3306} < ${SQL_DIR}/create_${DB_NAME}.sql\"
    expect \"Enter password:\"
    send \"${DB_PASS}\n\"
    interact
    "

# restore
expect -c "
    set timeout 30
    spawn sh -c \"mysql -u${DB_USER} -p -h${DB_HOST} ${DB_NAME} -P${DB_PORT:-3306} < ${SQL_DIR}/data_${DB_NAME}.sql\"
    expect \"Enter password:\"
    send \"${DB_PASS}\n\"
    interact
    "

# could be replaced like for `https://serverfault.com/questions/177135/expect-script-wait-command-hangs`
#expect -c "
#    set timeout 30
#    spawn sh -c \"mysql -u${DB_USER} -p -h${DB_HOST} ${DB_NAME} -P${DB_PORT:-3306} < ${SQL_DIR}/data_${DB_NAME}.sql\"
#    expect \"Enter password:\"
#    send \"${DB_PASS}\r\"
#    expect EOF
#    "
