#!/bin/sh
###############################################################################
# create test database
###############################################################################

# settings
DB_NAME=hiromaily2
DB_USER=root
DB_PASS=
DB_HOST=127.0.0.1
WORK_DIR=${GOPATH}/src/github.com/hiromaily/go-gin-wrapper/tests

# Create TestDB
expect -c "
    set timeout 30
    spawn sh -c \"mysql -u${DB_USER} -p${DB_PASS} -h${DB_HOST} < ${WORK_DIR}/createdb.sql\"
    expect \"Enter password:\"
    send \"${DB_PASS}\n\"
    interact
    "

# restore
expect -c "
    set timeout 30
    spawn sh -c \"mysql -u${DB_USER} -p${DB_PASS} -h${DB_HOST} ${DB_NAME} < ${WORK_DIR}/data.sql\"
    expect \"Enter password:\"
    send \"${DB_PASS}\n\"
    interact
    "
