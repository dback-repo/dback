@echo off
call check-command.cmd node --version || exit /b
::call check-command.cmd npm --help || exit /b
call check-command.cmd docker ps || exit /b
@echo on
