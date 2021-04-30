setlocal
@echo off

rem maven plugin version
set VER=%1
rem path to the golang source
set SRC=%2
set IMG_NAME=durable-task-binary-generator
docker build --build-arg PLATFORM=nanoserver -t %IMG_NAME%:%VER% .
docker run -i --rm --mount type=bind,src=%SRC%,dst=C:\durabletask %IMG_NAME%:%VER% C:\durabletask\test-and-compile.bat %VER%
docker rmi %IMG_NAME%:%VER%

@echo on
endlocal