cd %~dp0

docker rm -f dback-test1.1 dback-test1.2
docker volume rm dback-test1.1

docker volume create dback-test1.1

::infinity sleeping container
docker run -d --rm --name dback-test1.1 ^
	-v %CD%\data\mount-dir:/mount ^
	alpine:3.9.5 ^
	tail -f /dev/null 

::infinity sleeping container
docker run -d --rm --name dback-test1.2 ^
	-v %CD%\data\mount-dir:/mount ^
	alpine:3.9.5 ^
	tail -f /dev/null ::infinity sleeping container
