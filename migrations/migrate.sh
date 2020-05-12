#!/bin/sh

action=$1
mnumber=$2
db=$3

if [ "$action" = "composeup" ]; then
  /migrations/migrate \
    -source file:///migrations \
    -database "mysql://$VERIFICATION_DB_USER:$VERIFICATION_DB_PASSWORD@tcp($VERIFICATION_DB_HOST:$VERIFICATION_DB_PORT)/$VERIFICATION_DBNAME" \
    up
  exit 0
fi

if [ "$action" = "build" ]; then
  echo "building migrations"
  exit 0
fi

if [ "$action" != "goto" ] && [ "$action" != "force" ]; then
  echo "operation must be 'goto' or 'force'"
  exit 1
fi

if [ "$mnumber" = "" ]; then
  echo "a migration version is required"
  exit 1
fi

if [ "$db" = "" ]; then
  echo "a db address is required"
  exit 1
fi


/migrations/migrate \
    -source file:///migrations \
    -database $db \
    $action $mnumber
