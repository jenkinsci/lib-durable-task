setlocal
echo off

rem maven plugin version
set VER=%1
set BASEDIR=%CD%
set BIN=%BASEDIR%\bin
rem destination of the test reports
set TEST=%BASEDIR%\test-results
rem failure file
set FAILED=%TEST%\failed
set NAME=durable_task_monitor

rem gotestsum will generate junit test reports. v0.4.2 is the latest compatible with golang 1.14
del /s /q %TEST%
mkdir %TEST%
cd %BASEDIR%\pkg\common
go mod tidy
go get -v gotest.tools/gotestsum@v0.4.2
gotestsum --format standard-verbose --junitfile %TEST%\common-unit-tests.xml
if NOT %ERRORLEVEL% == 0 echo command>>%FAILED%
rem TODO test windows
rem cd %BASEDIR%\cmd\windows
rem go mod tidy
rem go get -v gotest.tools/gotestsum@v0.4.2
rem gotestsum --format standard-verbose --junitfile %TEST%\bash-unit-tests.xml
rem if NOT %ERRORLEVEL% == 0 echo windows>>%FAILED%

rem build the binaries
del /s /q %BIN%
mkdir %BIN%
cd %BASEDIR%/cmd/bash
set CGO_ENABLED=0& set GOOS=darwin& set GOARCH=amd64& go build -a -o %BIN%/%NAME%_%VER%_darwin_64
set CGO_ENABLED=0& set GOOS=darwin& set GOARCH=amd64& go build -a -o %BIN%/%NAME%_%VER%_darwin_32
set CGO_ENABLED=0& set GOOS=darwin& set GOARCH=amd64& go build -a -o %BIN%/%NAME%_%VER%_unix_64
set CGO_ENABLED=0& set GOOS=darwin& set GOARCH=amd64& go build -a -o %BIN%/%NAME%_%VER%_unix_32
rem TODO build windows
dir %BIN%

echo "binary generation complete. If you are still seeing this message, press Ctrl-C to exit."

echo on
endlocal