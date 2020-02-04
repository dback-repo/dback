cd %~dp0

docker rm -f dback-test1 dback-test2

::infinity sleeping container
docker run -d --rm --name dback-test1 ^
	alpine:3.9.5 ^
	tail -f /dev/null 

::infinity sleeping container
docker run -d --rm --name dback-test2 ^
	-v %CD%\test-data\mount:/mount ^
	alpine:3.9.5 ^
	tail -f /dev/null ::infinity sleeping container
