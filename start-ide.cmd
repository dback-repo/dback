@echo off
cd %~dp0
call scripts\setenv
cd %~dp0
call scripts\check-commands || exit /b
cd %~dp0


call npm i

echo ***************************************************
echo Type "npm run iter" for perform build and test
echo Change something in IDE, check the app can be compiled (ctrl+b), and run iter again. Repeat if required)
echo Find out all commands in "scripts" directory
echo ***************************************************

start liteide.exe %REPO%\go-app\src\dback\main.go

cmd
