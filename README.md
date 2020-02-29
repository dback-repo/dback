Dback is observe and backup all the mounts by pattern.<br>
By default, dback is also stop/start containers during backup, for prevent data corruption.
# The main use case
There is an instance with some important dockerized apps like Jenkins, GitLab, Vault, etc...
- Somebody wants to add an extra dockerized application, and he won't setup backup for the app
- All the apps can afford short downtime of each container
- All exist and extra dockerized apps of the instance must have periodical backup

# Example
```sh
mkdir /tmp/backup
docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v /tmp/backup:/backup dback:dback
```
Volume backups named as: `<ContainerName>/<PathInContainer>/snapshot-<timestamp>.tar`

##### Default pattern:
All mounts of each container matched all the options:
- HostConfig.AutoRemove: false
- HostConfig.RestartPolicy: always
- Status.State: running
- Status.Running: true

# Alternatives:
Perhaps these tools will work for you:<br>
https://github.com/christophetd/duplicacy-autobackup<br>
https://github.com/loomchild/volume-backup<br>
https://github.com/istepanov/docker-backup-to-s3<br>
https://github.com/lobaro/restic-backup-docker<br>