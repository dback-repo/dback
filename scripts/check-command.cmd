@echo off
%* > nul 2>&1 			:: run arguments as command, hide all output
if errorlevel 1 (
	@echo on
	echo "%*" has returned non-zero exit code. Check the application is installed and available.
	echo Press any key to exit.
	pause >nul
	exit /b %errorlevel%
)
@echo on
