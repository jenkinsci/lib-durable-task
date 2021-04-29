#! /bin/sh
set -ex
# maven plugin version
VER=$1
# path to the golang source
SRC=$2
IMG_NAME="durable-task-binary-generator"
BIN_NAME="durable_task_monitor"
docker build --build-arg PLATFORM="buster" \
             --build-arg ENTRY_SCRIPT="test-and-compile.sh" \
             --build-arg PLUGIN_VER=${VER} \
      -t ${IMG_NAME}:${VER} .
docker run -i --rm \
    --mount type=bind,src=${SRC},dst=/durabletask \
    ${IMG_NAME}:${VER}
docker rmi ${IMG_NAME}:${VER}