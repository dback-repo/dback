Dback is application for observe docker containers, make bulk incremental backups
of their mounts (folders and volumes), and pass backups to S3 bucket.

# The main use case
There is an instance with some important dockerized apps like Jenkins, GitLab, Vault, etc...
- Somebody wants to add an extra dockerized application, and he won't setup backup for the app
- All the apps can afford short downtime of each container
- All exist and extra dockerized apps of the instance must have periodical (daily) backup

# Example
Run test application. Zabbix.
```sh
docker run --name dback-example-zabbix --restart always -p 127.0.0.1:80:80 -d zabbix/zabbix-appliance:alpine-4.4.0
```
Run s3 server (minio), with a bucket
```sh
docker run --rm -d --name dback-test-1.minio -p 127.0.0.1:9000:9000 -e MINIO_ACCESS_KEY=dback_test -e MINIO_SECRET_KEY=3b464c70cf691ef6512ed51b2a minio/minio:RELEASE.2020-03-25T07-03-04Z server /data
docker run --rm -d --link dback-test-1.minio:minio --entrypoint=sh minio/mc:RELEASE.2020-05-28T23-43-36Z -c "mc config host add minio http://minio:9000 dback_test 3b464c70cf691ef6512ed51b2a && mc mb minio/dback-test"
```
Wait for ~1min for zabbix init<br>
Open http://localhost, check the login form shown<br>
<br>
Backup mounts of zabbix, and may be other containers match default selection pattern
```sh
docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback/dback:0.0.102 backup --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=SecureResticPassword11
```
Corrupt zabbix DB
```sh
docker exec -t dback-example-zabbix bash -c "rm -rf /var/lib/mysql/*"
```
Open http://localhost, check database error shown<br>
<br>
Restore all mounts of zabbix container
```sh
docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback/dback:0.0.102 restore container dback-example-zabbix --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=SecureResticPassword11
```
Open http://localhost, check login form works again

# Options
### Default containers selection pattern:
By default, backup will applied for all mounts of each container matched all the options:
- HostConfig.RestartPolicy == always
- HostConfig.AutoRemove == false
- Status.Running == true

You can override selection with --matcher flag. It is based on substrings matching in `docker inspect` json. It is awful, and planned to be updated with xpath matchers.


### Exclude mounts:
You able to ignore some mounts by regexp.<br>
`dback backup --exclude-mount "^/(drone.*|dback-test-1.5.*)$"`
this call will ignore all mounts started from "/drone" or "/dback-test-1.5"


# Alternatives:
Perhaps these tools will work for you:<br>
https://github.com/christophetd/duplicacy-autobackup<br>
https://github.com/loomchild/volume-backup<br>
https://github.com/istepanov/docker-backup-to-s3<br>
https://github.com/lobaro/restic-backup-docker<br>