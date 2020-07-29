const r = process.env.REPO
const cd = process.cwd()

const t = new(require(r+'/scripts/tools/tools.js'));
const ts = new(require(r+'/tests/shared/testing.js'));

ts.clearTmp()
try {cmd('docker rm -f dback-test-1.1 dback-test-1.2 dback-test-1.3 dback-test-1.4 dback-test-1.5 dback-test-1.6 dback-test-1.minio',{stdio: 'ignore'})}catch{}
try {cmd('docker volume rm dback-test-1.2-volume dback-test-1.2.1-volume dback-test-1.4-volume',{stdio: 'ignore'})}catch{}
ts.initVolumeWithFile('dback-test-1.2-volume','data/file1.txt')
ts.initVolumeWithFile('dback-test-1.2.1-volume','data/file1.txt')
ts.initVolumeWithFile('dback-test-1.4-volume','data/file1.txt')
ts.initVolumeWithFile('dback-test-1.6-volume','data/file1.txt')

//minio server (s3 compatible) for saving test mounts
t.cmd('docker run --rm -d --name dback-test-1.minio -p 127.0.0.1:2157:9000 -e MINIO_ACCESS_KEY=dback_test -e MINIO_SECRET_KEY=3b464c70cf691ef6512ed51b2a minio/minio:RELEASE.2020-03-25T07-03-04Z server /data')
t.cmd('docker run --rm -d --link dback-test-1.minio:minio --entrypoint=sh minio/mc:RELEASE.2020-05-28T23-43-36Z -c "mc config host add minio http://minio:9000 dback_test 3b464c70cf691ef6512ed51b2a && mc mb minio/dback-test"')

//this containers should be saved
t.cmd('docker run --restart always -d --name dback-test-1.1 -v '+cd+'/data/mount-dir:/mount-dir nginx:1.17.8-alpine')
t.cmd('docker run --link dback-test-1.1:dback-test-1.1 --restart always -d --name dback-test-1.2 -v '+cd+'/data/mount-dir:/mount-dir -v dback-test-1.2-volume:/mount-vol -v dback-test-1.2.1-volume:/mount-vol-for-exclude nginx:1.17.8-alpine')

//this containers should be ignored
t.cmd('docker run --link dback-test-1.2:dback-test-1.2 --restart always -d --name dback-test-1.3 nginx:1.17.8-alpine') 														//container has no mounts
t.cmd('docker run --link dback-test-1.3:dback-test-1.3 --rm -d --name dback-test-1.4 -v dback-test-1.4-volume:/mount-vol nginx:1.17.8-alpine') 				//temporary container (--rm)
t.cmd('docker run --link dback-test-1.4:dback-test-1.4 --restart always -d --name dback-test-1.5 -v '+cd+'/data/mount-dir:/mount-dir nginx:1.17.8-alpine') 					//ignored by --exclude-mount pattern
t.cmd('docker run --link dback-test-1.5:dback-test-1.5 -d --name dback-test-1.6 -v '+cd+'/data/mount-dir:/mount-dir -v dback-test-1.6-volume:/mount-vol nginx:1.17.8-alpine')  //ignored due restart-policy==none


// var out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback backup -e -x "^/(drone.*|dback-test-1.5.*)$" -x "for-exclude$"').toString()
// ts.checkSub(out,'Ignore container:  /dback-test-1.6  cause: matcher not found 'RestartPolicy':{"Name":"always"')
// ts.checkSub(out,'Ignore container:  /dback-test-1.4  cause: matcher not found 'RestartPolicy':{"Name":"always"')
// ts.checkSub(out,'Ignore container:  /dback-test-1.3  cause: container has no mounts')
// ts.checkSub(out,'Ignore container:  /dback-test-1.minio  cause: matcher not found 'RestartPolicy':{"Name":"always"')
// ts.checkSub(out,'Exclude mount: /dback-test-1.5/mount-dir      cause: --exclude-mount ^/(drone.*|dback-test-1.5.*)$')
// ts.checkSub(out,'Exclude mount: /dback-test-1.2/mount-vol-for-exclude      cause: --exclude-mount for-exclude$')
// ts.checkSub(out,'Emulation started')
// ts.checkSub(out,'/dback-test-1.1/mount-dir')
// ts.checkSub(out,'/dback-test-1.2/mount-dir')
// ts.checkSub(out,'/dback-test-1.2/mount-vol')
// ts.checkSub(out,'The mounts above will be backup, if run dback without --emulate (-e) flag')

var out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback backup -x "^/(drone.*|dback-test-1.5.*)$" -x "for-exclude$" --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
ts.checkSub(out,'Ignore container:  /dback-test-1.6  cause: matcher not found //*[@name=\'RestartPolicy\']/string[text()!=\'no\']')
ts.checkSub(out,'Ignore container:  /dback-test-1.4  cause: matcher not found //*[@name=\'RestartPolicy\']/string[text()!=\'no\']')
ts.checkSub(out,'Ignore container:  /dback-test-1.3  cause: container has no mounts')
ts.checkSub(out,'Ignore container:  /dback-test-1.minio  cause: matcher not found //*[@name=\'RestartPolicy\']/string[text()!=\'no\']')
ts.checkSub(out,'Exclude mount: /dback-test-1.5/mount-dir      cause: --exclude-mount ^/(drone.*|dback-test-1.5.*)$')
ts.checkSub(out,'Exclude mount: /dback-test-1.2/mount-vol-for-exclude      cause: --exclude-mount for-exclude$')
ts.checkSub(out,'Backup started. Timestamp = ')
ts.checkSub(out,'Save to restic: /dback-test-1.1/mount-dir')
ts.checkSub(out,'Save to restic: /dback-test-1.2/mount-dir')
ts.checkSub(out,'Save to restic: /dback-test-1.2/mount-vol')
ts.checkSub(out,'Backup finished for the mounts above, in ')
console.log(out)

//make 2nd snapshot, ignore container dback-test-1.2
var out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback backup -x "^/(drone.*|dback-test-1.5.*|dback-test-1.2.*)$" -x "for-exclude$" --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()

//re-create dback-test-1.2, with new volumes 
t.cmd('docker rm -f dback-test-1.2')
t.cmd('docker volume rm dback-test-1.2-volume')
//t.cmd('docker volume create dback-test-1.2-volume')
ts.initVolumeWithFile('dback-test-1.2-volume','data/file3.txt')
t.cmd('docker run --restart always -d --name dback-test-1.2 -v dback-test-1.2-volume:/mount-vol nginx:1.17.8-alpine')

//restore 
out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback restore --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
console.log(out)

// //check restored volume
t.cmd('docker cp dback-test-1.2:/mount-vol '+cd+'/tmp')
// var content=fs.readFileSync(cd+'/tmp/mount-vol/file1.txt', "utf8");
// if (content!='file1'){
// 	throw('File content is invalid. "file1" expected, but actually is "'+content+'"')
// }

//docker exec -t dback-test-1.2 sh -c cat /mount-vol/1.txt



//var out = t.cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v '+cd+'/tmp:/dback-snapshots dback backup --exclude-mount "^/(drone.*|dback-test-1.5.*)$" '+process.env.S3_ENDPOINT+' '+process.env.S3_BUCKET+' '+process.env.ACC_KEY+' '+process.env.SEC_KEY).toString()

// const t = new(require(r+'/scripts/tools/tools.js'))
// out = t.t.cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v '+cd+'/tmp:/dback-snapshots dback restore '+process.env.S3_ENDPOINT+' '+process.env.S3_BUCKET+' '+process.env.ACC_KEY+' '+process.env.SEC_KEY).toString()
// console.log(out)


//restore single container with same name
out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback restore container dback-test-1.1 --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
console.log(out)

//restore single mount with same container name name
out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback restore mount /dback-test-1.1/mount-dir --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
console.log(out)

//restore single container with diff name
out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback restore container dback-test-1.1 dback-test-1.2 --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
console.log(out)

//restore single mount with diff container name
out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback restore mount /dback-test-1.1/mount-dir /dback-test-1.2/mount-dir --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
console.log(out)

//list mounts
out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback ls --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
console.log(out)

//list mounts with prefix
out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback ls /dback-test --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
console.log(out)

//list snapshots of the mount
out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback ls /dback-test-1.1/mount-dir --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
console.log(out)

//list snapshots with prefix doesn't exist
out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback ls /containerNotExist --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
console.log(out)

//show xml of container inspect
out = t.cmd('docker run --rm -t -v //var/run/docker.sock:/var/run/docker.sock dback inspect /dback-test-1.1').toString()
console.log(out)


//tests for restore all the kinds, with snapshot

//restore single container with same name, with incorrect snapshot
// out = t.cmd('docker run --rm -t --link dback-test-1.minio:minio -v //var/run/docker.sock:/var/run/docker.sock dback restore container dback-test-1.1 --snapshot=09.07.2020.16-30-25 --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
// console.log(out)




t.cmd('docker rm -f dback-test-1.1 dback-test-1.2 dback-test-1.3 dback-test-1.4 dback-test-1.5 dback-test-1.6 dback-test-1.minio')
t.cmd('docker volume rm dback-test-1.2-volume dback-test-1.2.1-volume dback-test-1.4-volume dback-test-1.6-volume') 