@echo off
setlocal

rem docker build --output requires BuildKit, not available for windows containers
rem see https://github.com/microsoft/Windows-Containers/issues/34
rem instead, create a temporary writeable container layer to copy out the binaries
set VER=%1
set DEST=%2
set IMG_NAME=durable-task-binary-generator
set BINARY_NAME=durable_task_monitor
set OUTPUT_DIR=/durabletask/cmd/bash
docker build --build-arg VERSION=%VER% -f Dockerfile.windows -t %IMG_NAME%:%VER% .
docker create -ti --name scratch %IMG_NAME%:%VER%
docker cp scratch:%OUTPUT_DIR%/%BINARY_NAME%_%VER%_darwin_amd_64 %DEST%
docker cp scratch:%OUTPUT_DIR%/%BINARY_NAME%_%VER%_darwin_arm_64 %DEST%
docker cp scratch:%OUTPUT_DIR%/%BINARY_NAME%_%VER%_linux_64 %DEST%
docker cp scratch:%OUTPUT_DIR%/%BINARY_NAME%_%VER%_linux_32 %DEST%
docker rm -f scratch
docker rmi %IMG_NAME%:%VER%

endlocal
@echo on