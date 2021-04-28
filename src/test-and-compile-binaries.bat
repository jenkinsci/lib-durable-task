rem TODO
setlocal

rem maven plugin version
set VER=%1
set BASEDIR=%CD%
set BIN=%BASEDIR%\bin
rem destination of the test reports
set TEST=%BASEDIR%\test-results
rem failure file
set FAILED=%TEST%\failed
set NAME=durable_task_monitor

echo %VER%
echo %BASEDIR%
echo %BIN%
echo %TEST%
echo %FAILED%
echo %NAME%

endlocal