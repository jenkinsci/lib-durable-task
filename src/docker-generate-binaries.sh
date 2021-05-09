#! /bin/sh
set -ex
# maven plugin version
VER=$1
# output directory of binaries
DEST=$2
docker build --build-arg VERSION=$VER -o $DEST -f Dockerfile.linux .