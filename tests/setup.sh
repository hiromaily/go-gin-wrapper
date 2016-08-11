#!/bin/sh
###############################################################################
# create test database
###############################################################################

# settings
DB_HOST=127.0.0.1
DB_NAME=hiromaily2
DB_USER=root
#DB_PASS=
#DB_PORT=3306
WORK_DIR=${GOPATH}/src/github.com/hiromaily/go-gin-wrapper/tests


# Create TestDB
expect -c "
    set timeout 30
    spawn sh -c \"mysql -u${DB_USER} -p -h${DB_HOST} -P${DB_PORT:-3306} < ${WORK_DIR}/createdb.sql\"
    expect \"Enter password:\"
    send \"${DB_PASS}\n\"
    interact
    "

# restore
expect -c "
    set timeout 30
    spawn sh -c \"mysql -u${DB_USER} -p -h${DB_HOST} ${DB_NAME} -P${DB_PORT:-3306} < ${WORK_DIR}/data.sql\"
    expect \"Enter password:\"
    send \"${DB_PASS}\n\"
    interact
    "
