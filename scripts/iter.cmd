cd %~dp0
call build.cmd  || exit /b
call test.cmd   || exit /b