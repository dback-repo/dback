const mv = require('fs').renameSync
var fs  = require('fs');
const cmd = require('child_process').execSync
const r = process.env.REPO
const cd = process.cwd()
const Path = require('path');
const t = new(require(r+'/scripts/tools/tools.js'));


const createFolder = function(dir) {
	if (!fs.existsSync(dir)){
		fs.mkdirSync(dir);
	}
};

//one does not simply copy file to the volume!
//we must mount the volume to temporary container
const moveFileToDockerVolume = function(file, volume){
	t.cmd('docker run --rm -v '+Path.resolve(file)+':/'+Path.basename(file)+' -v '+volume+':/dest busybox cp -r /'+Path.basename(file)+' /dest')
}

//move an old tmp folder to tmp-old
//because of deleting is extrimely slow for multiple files
const clearTmp = function(){
	createFolder('tmp')
	createFolder('tmp-old')
	fs.renameSync('tmp', 'tmp-old/'+Math.floor(new Date() / 1000))
	createFolder('tmp')
}

const initVolumeWithFile = function(volumeName, filename){
	try {t.cmd('docker volume rm ' + volumeName)}catch{}
	t.cmd('docker volume create ' + volumeName)
	moveFileToDockerVolume(filename, volumeName)
}

const checkSub = function(str, substr){
	if (!str.includes(substr)){
		throw('"'+substr+'" is expected but not found in the string:'+"\r\n"+str)
	}
}

const checkNoSub = function(str, substr){
	if (str.includes(substr)){
		throw('"'+substr+'" is expected but not found in the string:'+"\r\n"+str)
	}
}


clearTmp()
try {cmd('docker rm -f dback-test-1.1 dback-test-1.2 dback-test-1.3 dback-test-1.4 dback-test-1.5 dback-test-1.6 dback-test-minio',{stdio: 'ignore'})}catch{}
initVolumeWithFile('dback-test-1.2-volume','data/file1.txt')
initVolumeWithFile('dback-test-1.4-volume','data/file1.txt')

//minio server (s3 compatible) for saving test mounts
t.cmd('docker run --rm -d --name dback-test-minio -p 127.0.0.1:2157:9000 -e MINIO_ACCESS_KEY=dback_test -e MINIO_SECRET_KEY=3b464c70cf691ef6512ed51b2a minio/minio:RELEASE.2020-03-25T07-03-04Z server /data')
t.cmd('docker run --rm -d --link dback-test-minio:minio --entrypoint=sh minio/mc:RELEASE.2020-05-28T23-43-36Z -c "mc config host add minio http://minio:9000 dback_test 3b464c70cf691ef6512ed51b2a && mc mb minio/dback-test"')

//this containers should be saved
t.cmd('docker run --restart always -d --name dback-test-1.1 -v '+cd+'/data/mount-dir:/mount-dir nginx:1.17.8-alpine')
t.cmd('docker run --restart on-failure -d --name dback-test-1.2 -v '+cd+'/data/mount-dir:/mount-dir -v dback-test-1.2-volume:/mount-vol nginx:1.17.8-alpine')

//this containers should be ignored
t.cmd('docker run --restart always -d --name dback-test-1.3 nginx:1.17.8-alpine') 														//container has no mounts
t.cmd('docker run --rm -d --name dback-test-1.4 -v dback-test-1.4-volume:/mount-vol nginx:1.17.8-alpine') 				//temporary container (--rm)
t.cmd('docker run --restart always -d --name dback-test-1.5 -v '+cd+'/data/mount-dir:/mount-dir nginx:1.17.8-alpine') 					//ignored by --exclude-mount pattern
t.cmd('docker run -d --name dback-test-1.6 -v '+cd+'/data/mount-dir:/mount-dir -v dback-test-1.2-volume:/mount-vol nginx:1.17.8-alpine')  //ignored due restart-policy==none

var out = t.cmd('docker run --rm -t --link dback-test-minio:minio -v //var/run/docker.sock:/var/run/docker.sock -v '+cd+'/tmp:/dback-data dback backup -e -x "^/(drone.*|dback-test-1.5.*)$" --s3-endpoint=http://minio:9000 -b=dback-test -a=dback_test -s=3b464c70cf691ef6512ed51b2a -p=sdf').toString()
// checkSub(out,'Backup started')
// checkSub(out,'exclude: /dback-test-1.4      Reason: temporary container (--rm)')
// checkSub(out,'exclude: /dback-test-1.5/mount-dir      Reason: --exclude-mount parameter')
// checkSub(out,'exclude: /dback-test-1.6      Reason: container restart policy==none')
// checkSub(out,'make backup: /dback-test-1.2/mount-vol')
// checkSub(out,'make backup: /dback-test-1.1/mount-dir')
// checkSub(out,'make backup: /dback-test-1.2/mount-dir')
// checkSub(out,'exclude: /dback-test-1.5/mount-dir')
// checkSub(out,'Backup has finished for the mounts above')

console.log(out)

//var out = t.cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v '+cd+'/tmp:/dback-snapshots dback backup --exclude-mount "^/(drone.*|dback-test-1.5.*)$" '+process.env.S3_ENDPOINT+' '+process.env.S3_BUCKET+' '+process.env.ACC_KEY+' '+process.env.SEC_KEY).toString()

// const t = new(require(r+'/scripts/tools/tools.js'))
// out = t.t.cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v '+cd+'/tmp:/dback-snapshots dback restore '+process.env.S3_ENDPOINT+' '+process.env.S3_BUCKET+' '+process.env.ACC_KEY+' '+process.env.SEC_KEY).toString()
// console.log(out)

t.cmd('docker rm -f dback-test-1.1 dback-test-1.2 dback-test-1.3 dback-test-1.4 dback-test-1.5 dback-test-1.6 dback-test-minio')