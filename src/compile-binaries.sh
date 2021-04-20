#! /bin/sh
set -x
# maven plugin version
VER=$1
NAME="durable_task_monitor"
rm -rf ${NAME}_*
cd "cmd/bash"
go mod tidy
env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o ../../${NAME}_${VER}_darwin_64
env CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -a -o ../../${NAME}_${VER}_darwin_32
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ../../${NAME}_${VER}_unix_64
env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -o ../../${NAME}_${VER}_unix_32
#TODO: Windows binary generation
#cd "../windows"
#go mod tidy
#env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o ../../${NAME}_${VER}_win_64
#env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -a -o ../../${NAME}_${VER}_win_32
echo "binary generation complete. If you are still seeing this message, press Ctrl-C to exit."
