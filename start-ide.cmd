@echo off
cd %~dp0
call scripts\setenv
cd %~dp0
call scripts\check-commands || exit /b
cd %~dp0


call npm i

echo ***************************************************
echo Type "iter" for build and test, change something, and run tests again. Repeat if required)
echo Find out all commands in "scripts" directory
echo ***************************************************

start liteide.exe %REPO%\go-app\src\dback\main.go

cmd
