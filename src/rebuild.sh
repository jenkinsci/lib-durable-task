#! /bin/sh
# Convenience script to rebuild golang binaries during development
if [[ $1 -eq 0 ]] ; then
    echo 'please provide a plugin version as an argument (ex: 1.32)'
    exit 0
fi
set -x
# maven plugin version
VER=$1
NAME="durable_task_monitor"
rm -rf ${NAME}_*
cd "cmd/bash"
env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o ../../${NAME}_${VER}_darwin_64
env CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -a -o ../../${NAME}_${VER}_darwin_32
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ../../${NAME}_${VER}_unix_64
env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -o ../../${NAME}_${VER}_unix_32
#TODO: Windows
#cd "../windows"
#env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o ../../${NAME}_${VER}_win_64
#env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -a -o ../../${NAME}_${VER}_win_32
