#!/bin/bash
set -eux
docker build \
    -t docker.pkg.github.com/ninja-software/infinote/infinote \
    --no-cache \
    --network=host \
    --build-arg GOPROXY_DEFAULT=http://localhost:3000 \
    --build-arg FONTAWESOME_TOKEN_DEFAULT=$FONTAWESOME_TOKEN \
    .

docker push docker.pkg.github.com/ninja-software/infinote/infinote