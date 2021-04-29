setlocal

rem maven plugin version
set VER=%1
rem path to the golang source
set SRC=%2
set IMG_NAME=durable-task-binary-generator
set BIN_NAME=durable_task_monitor
docker build --build-arg PLATFORM=nanoserver ^
             --build-arg ENTRY_SCRIPT=test-and-compile.bat ^
             --build-arg PLUGIN_VER=%VER% ^
      -t %IMG_NAME%:%VER% .
docker run -i --rm ^
    --mount type=bind,src=%SRC%,dst=C:\durabletask ^
    %IMG_NAME%:%VER%
rem docker rmi %IMG_NAME%:%VER%

endlocal