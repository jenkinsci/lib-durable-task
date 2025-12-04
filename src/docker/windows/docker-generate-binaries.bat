@echo off
setlocal

rem docker build --output requires BuildKit, not available for windows containers
rem see https://github.com/microsoft/Windows-Containers/issues/34
rem instead, create a temporary writeable container layer to copy out the binaries

rem output directory of binaries
set DEST=%1
set IMG_NAME=durable-task-binary-generator
set BINARY_NAME=durable_task_monitor
set BASH_DIR=/durabletask/cmd/bash
set WIN_DIR=/durabletask/cmd/windows
mkdir "%DEST%"
docker build --no-cache -t %IMG_NAME%:0.0 -f docker/windows/Dockerfile .
docker create -ti --name scratch %IMG_NAME%:0.0
docker cp scratch:%BASH_DIR%/%BINARY_NAME%_darwin_amd64 %DEST%
docker cp scratch:%BASH_DIR%/%BINARY_NAME%_darwin_arm64 %DEST%
docker cp scratch:%BASH_DIR%/%BINARY_NAME%_linux_64 %DEST%
docker cp scratch:%BASH_DIR%/%BINARY_NAME%_linux_32 %DEST%
docker cp scratch:%BASH_DIR%/%BINARY_NAME%_linux_ppc64le %DEST%
docker cp scratch:%BASH_DIR%/%BINARY_NAME%_linux_aarch64 %DEST%
docker cp scratch:%WIN_DIR%/%BINARY_NAME%_win_64.exe %DEST%
docker cp scratch:%WIN_DIR%/%BINARY_NAME%_win_32.exe %DEST%
docker rm -f scratch
docker rmi %IMG_NAME%:0.0

endlocal
@echo on
