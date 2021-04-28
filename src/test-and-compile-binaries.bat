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

rem gotestsum will generate junit test reports. v0.4.2 is the latest compatible with golang 1.14
del %TEST%
mkdir %TEST%
cd %BASEDIR5\pkg\common
go mod tidy
go get -v gotest.tools/gotestsum@v0.4.2
gotestsum --format standard-verbose --junitfile %TEST%\common-unit-tests.xml
if %ERROR_LEVEL%  == 0 echo command unit tests success
cd %BASEDIR%\cmd\bash
go mod tidy
go get -v gotest.tools/gotestsum@v0.4.2
gotestsum --format standard-verbose --junitfile %TEST%\bash-unit-tests.xml
if %ERROR_LEVEL%  == 0 echo bash unit tests success

endlocal