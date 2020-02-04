cd %~dp0

docker rm -f dback-test1.1 dback-test1.2
docker volume rm dback-test1.1

docker volume create dback-test1.1  || exit /b

::infinity sleeping container
docker run -d --rm --name dback-test1.1 ^
	-v %CD%\data\mount-dir:/mount ^
	docker:19.03.5 ^
	tail -f /dev/null  || exit /b

::infinity sleeping container
docker run -d --rm --name dback-test1.2 ^
	-v %CD%\data\mount-dir:/mount ^
	docker:19.03.5 ^
	tail -f /dev/null  || exit /b