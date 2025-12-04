#! /bin/sh
set -ex
# output directory of binaries
DEST=$1
export DOCKER_BUILDKIT=1
docker build --no-cache -o $DEST -f docker/linux/Dockerfile .
