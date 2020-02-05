cd %REPO%\go-app\src\dback
set GOOS=linux
call go build || exit /b
ren %REPO%\go-app\src\dback\dback %REPO%\docker\dback || exit /b

cd %REPO%\docker
docker build -t dback .  || exit /b
