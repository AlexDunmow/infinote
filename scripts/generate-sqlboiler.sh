#!/bin/bash
REPO=`git rev-parse --show-toplevel`
$REPO/bin/sqlboiler $REPO/bin/sqlboiler-psql --wipe --tag db --config $REPO/server/sqlboiler.toml --output $REPO/server/db