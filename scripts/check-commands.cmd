node --version > nul 2>&1
if errorlevel 1 (
	echo "node --version" has returned non-zero exit code. Check the application is installed and available.
	echo Press any key to exit.
	pause >nul
	exit /b %errorlevel%
)

call npm --version > nul 2>&1
if errorlevel 1 (
	echo "npm" has returned non-zero exit code. Check the application is installed and available.
	echo Press any key to exit.
	pause >nul
	exit /b %errorlevel%
)

docker ps > nul 2>&1
if errorlevel 1 (
	echo "docker ps" has returned non-zero exit code. Check the application is installed and available.
	echo Press any key to exit.
	pause >nul
	exit /b %errorlevel%
)