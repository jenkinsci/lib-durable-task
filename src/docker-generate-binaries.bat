@echo off
setlocal

rem docker build --output requires BuildKit, not available for windows containers
rem see https://github.com/microsoft/Windows-Containers/issues/34
rem instead, create a temporary writeable container layer to copy out the binaries

rem output directory of binaries
set DEST=%1
set IMG_NAME=durable-task-binary-generator
set BINARY_NAME=durable_task_monitor
set OUTPUT_DIR=/durabletask/cmd/bash
mkdir "%DEST%"
docker build --no-cache -f Dockerfile.windows -t %IMG_NAME%:0.0 .
docker create -ti --name scratch %IMG_NAME%:0.0
docker cp scratch:%OUTPUT_DIR%/%BINARY_NAME%_darwin_amd64 %DEST%
docker cp scratch:%OUTPUT_DIR%/%BINARY_NAME%_darwin_arm64 %DEST%
docker cp scratch:%OUTPUT_DIR%/%BINARY_NAME%_linux_64 %DEST%
docker cp scratch:%OUTPUT_DIR%/%BINARY_NAME%_linux_32 %DEST%
docker rm -f scratch
docker rmi %IMG_NAME%:0.0

endlocal
@echo on