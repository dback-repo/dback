const mv = require('fs').renameSync
var fs  = require('fs');
const cmd = require('child_process').execSync
const spw = require('child_process').spawnSync
const r = process.env.REPO
const cd = process.cwd()

const Path = require('path');

const createFolder = function(dir) {
	if (!fs.existsSync(dir)){
		fs.mkdirSync(dir);
	}
};

//one does not simply copy file to the volume!
//we must mount the volume to temporary container
const moveFileToDockerVolume = function(file, volume){
	cmd('docker run --rm -v '+Path.resolve(file)+':/'+Path.basename(file)+' -v '+volume+':/dest busybox cp -r /'+Path.basename(file)+' /dest')
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
	try {cmd('docker volume rm ' + volumeName)}catch{}
	cmd('docker volume create ' + volumeName)
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
try {cmd('docker rm -f dback-test-1.1 dback-test-1.2 dback-test-1.3 dback-test-1.4 dback-test-1.5 dback-test-1.6',{stdio: 'ignore'})}catch{}
initVolumeWithFile('dback-test-1.2-volume','data/file1.txt')
initVolumeWithFile('dback-test-1.4-volume','data/file1.txt')

//this containers should be saved
cmd('docker run --restart always -d --name dback-test-1.1 -v '+cd+'/data/mount-dir:/mount-dir nginx:1.17.8-alpine')
cmd('docker run --restart on-failure -d --name dback-test-1.2 -v '+cd+'/data/mount-dir:/mount-dir -v dback-test-1.2-volume:/mount-vol nginx:1.17.8-alpine')

//this containers should be ignored
cmd('docker run --restart always -d --name dback-test-1.3 nginx:1.17.8-alpine') 														//container has no mounts
cmd('docker run --rm -d --name dback-test-1.4 -v dback-test-1.4-volume:/mount-vol nginx:1.17.8-alpine') 				//temporary container (--rm)
cmd('docker run --restart always -d --name dback-test-1.5 -v '+cd+'/data/mount-dir:/mount-dir nginx:1.17.8-alpine') 					//ignored by --exclude-mount pattern
cmd('docker run -d --name dback-test-1.6 -v '+cd+'/data/mount-dir:/mount-dir -v dback-test-1.2-volume:/mount-vol nginx:1.17.8-alpine')  //ignored due restart-policy==none

var out = cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v '+cd+'/tmp:/dback-snapshots dback backup --exclude-mount "^/(drone.*|dback-test-1.5.*)$" '+process.env.S3_ENDPOINT+' '+process.env.S3_BUCKET+' '+process.env.ACC_KEY+' '+process.env.SEC_KEY).toString()
checkSub(out,'Backup started')
checkSub(out,'exclude: /dback-test-1.4      Reason: temporary container (--rm)')
checkSub(out,'exclude: /dback-test-1.5/mount-dir      Reason: --exclude-mount parameter')
checkSub(out,'exclude: /dback-test-1.6      Reason: container restart policy==none')
checkSub(out,'make backup: /dback-test-1.2/mount-vol')
checkSub(out,'make backup: /dback-test-1.1/mount-dir')
checkSub(out,'make backup: /dback-test-1.2/mount-dir')
checkSub(out,'exclude: /dback-test-1.5/mount-dir')
checkSub(out,'Backup has finished for the mounts above')

console.log(out)

//var out = cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v '+cd+'/tmp:/dback-snapshots dback backup --exclude-mount "^/(drone.*|dback-test-1.5.*)$" '+process.env.S3_ENDPOINT+' '+process.env.S3_BUCKET+' '+process.env.ACC_KEY+' '+process.env.SEC_KEY).toString()

const t = new(require(r+'/scripts/tools/tools.js'))
out = t.cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v '+cd+'/tmp:/dback-snapshots dback restore '+process.env.S3_ENDPOINT+' '+process.env.S3_BUCKET+' '+process.env.ACC_KEY+' '+process.env.SEC_KEY).toString()

console.log(out)

cmd('docker rm -f dback-test-1.1 dback-test-1.2 dback-test-1.3 dback-test-1.4 dback-test-1.5 dback-test-1.6')