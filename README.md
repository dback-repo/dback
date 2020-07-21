Dback is application for observe docker containers, make bulk incremental backups
of their mounts (folders and volumes), and pass backups to S3 bucket. Restoring to exist containers is also available.<br>
Runs [restic](https://github.com/restic/restic) under the hood.

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
Backup mounts of zabbix, and may be some other containers matching default selection pattern
```sh
docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback/dback:0.0.110 backup --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=SecureResticPassword11
```
Corrupt zabbix DB
```sh
docker exec -t dback-example-zabbix bash -c "rm -rf /var/lib/mysql/*"
```
Open http://localhost, check database error shown<br>
<br>
Restore all mounts of zabbix container
```sh
docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback/dback:0.0.110 restore container dback-example-zabbix --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=SecureResticPassword11
```
Open http://localhost, check login form works again

# Backup options
### Default containers selection pattern:
By default, backup will be applied for all mounts of each container matched the options:
- HostConfig.RestartPolicy != no
- HostConfig.AutoRemove == false
- Status.Running == true

The pattern can be overridden with --matcher flag. It is based on xpath matching in `docker inspect` xml. You can see the xml with `dback inspect <container>`.


### Exclude mounts or containers:
You able to ignore some mounts by regex.<br>
`dback backup --exclude-mount "^/(drone.*|dback-test-1.5.*)$"`
this call will ignore all mounts started from "/drone" or "/dback-test-1.5"


# Alternatives, and why dback created:
Inspired by these projects:<br>
https://github.com/blacklabelops/volumerize<br>
https://github.com/christophetd/duplicacy-autobackup<br>
https://github.com/loomchild/volume-backup<br>
https://github.com/istepanov/docker-backup-to-s3<br>
https://github.com/lobaro/restic-backup-docker<br>
<br>
The apps listed above did not support some features dback provides:
* **Containers observation.** Dback is find and backup all containers matching the pattern.
* **Auto stop/start containers.** Dback is always stops containers before make backup, and then start they, even if something went wrong.
* **Incremental backup.** 2nd and subsequent backups are faster than 1st, because [restic](https://github.com/restic/restic) sending and store only the difference between shots.
* **Bulk restoring.** You are able to restore all the mounts of target container, and even all saved mounts of all containers, with a single command.
* **Concurrency.** By default - start a thread for each mount. Can be decreased.
* **Retrying.**  [Restic](https://github.com/restic/restic) will retry mount backup procedure, if something went wrong.

That's why dback is more useful in some cases.