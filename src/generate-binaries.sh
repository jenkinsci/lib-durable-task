#! /bin/sh
set -ex
# maven plugin version
VER=$1
# path to the golang source
SRC=$2
IMG_NAME="durable-task-binary-generator"
docker build --build-arg PLATFORM="buster" -t ${IMG_NAME}:${VER} .
docker run -i --rm --mount type=bind,src=${SRC},dst=/durabletask ${IMG_NAME}:${VER} /durabletask/test-and-compile.sh ${VER}
docker rmi ${IMG_NAME}:${VER}