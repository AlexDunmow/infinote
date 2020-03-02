#!/bin/bash

REPO=`git rev-parse --show-toplevel`

source $REPO/scripts/db-prepare-userpass.sh

$REPO/bin/migrate -database "postgres://$LOCAL_DEV_DB_USER:$LOCAL_DEV_DB_PASS@$LOCAL_DEV_DB_HOST:$LOCAL_DEV_DB_PORT/$LOCAL_DEV_DB_DATABASE?sslmode=disable" -path ./migrations drop
$REPO/bin/migrate -database "postgres://$LOCAL_DEV_DB_USER:$LOCAL_DEV_DB_PASS@$LOCAL_DEV_DB_HOST:$LOCAL_DEV_DB_PORT/$LOCAL_DEV_DB_DATABASE?sslmode=disable" -path ./migrations up
