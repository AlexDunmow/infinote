#!/bin/bash
set -e
source /var/scripts/set-environment.sh

echo "Pulling container..."

echo $DOCKER_TOKEN | docker login docker.pkg.github.com -u nii236 --password-stdin &> /dev/null
docker pull --quiet  docker.pkg.github.com/ninja-software/infinote/infinote:latest &> /dev/null

echo "Removing existing containers..."

[ "$(docker ps -a | grep infinote)" ] && docker stop infinote &> /dev/null
[ "$(docker ps -a | grep infinote)" ] && docker rm infinote &> /dev/null

echo "Migrating database..."

docker run \
        --rm \
        --name infinote-db-client \
        -e BOILERPLATE_LOADBALANCER_ROOTPATH=$BOILERPLATE_LOADBALANCER_ROOTPATH \
        -e BOILERPLATE_DATABASE_USER=$BOILERPLATE_DATABASE_USER \
        -e BOILERPLATE_DATABASE_PASS=$BOILERPLATE_DATABASE_PASS \
        -e BOILERPLATE_DATABASE_HOST=$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+') \
        -e BOILERPLATE_DATABASE_PORT=$BOILERPLATE_DATABASE_PORT \
        -e BOILERPLATE_DATABASE_NAME=$BOILERPLATE_DATABASE_NAME \
docker.pkg.github.com/ninja-software/infinote/infinote:latest db-migrate

echo "Running server..."

docker run \
        --detach \
        --name infinote \
        -p 8082:8080 \
        -e BOILERPLATE_LOADBALANCER_ROOTPATH=$BOILERPLATE_LOADBALANCER_ROOTPATH \
        -e BOILERPLATE_DATABASE_USER=$BOILERPLATE_DATABASE_USER \
        -e BOILERPLATE_DATABASE_PASS=$BOILERPLATE_DATABASE_PASS \
        -e BOILERPLATE_DATABASE_HOST=$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+') \
        -e BOILERPLATE_DATABASE_PORT=$BOILERPLATE_DATABASE_PORT \
        -e BOILERPLATE_DATABASE_NAME=$BOILERPLATE_DATABASE_NAME \
docker.pkg.github.com/ninja-software/infinote/infinote:latest
