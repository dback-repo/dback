const mv = require('fs').renameSync
var fs  = require('fs');
const cmd = require('child_process').execSync
const r = process.env.REPO

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
	createFolder('tmp-old')
	fs.renameSync('tmp', 'tmp-old/'+Math.floor(new Date() / 1000))
	createFolder('tmp')
}

const initVolumeWithFile = function(){
	try {cmd('docker volume rm dback-test-1.2-volume')}catch{}
	cmd('docker volume create dback-test-1.2-volume')
	moveFileToDockerVolume('data/file1.txt','dback-test-1.2-volume')
}


clearTmp()
try {cmd('docker rm -f dback-test-1.1 dback-test-1.2 dback-test-1.3')}catch{}
initVolumeWithFile()
cmd('docker run -d --name dback-test-1.1 -v %CD%\\data\\mount-dir:/mount-dir nginx:1.17.8-alpine')
cmd('docker run -d --name dback-test-1.2 -v %CD%\\data\\mount-dir:/mount-dir -v dback-test-1.2-volume:/mount-vol nginx:1.17.8-alpine')
cmd('docker run -d --name dback-test-1.3 nginx:1.17.8-alpine')
cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v %CD%\\tmp:/backup dback',{stdio: 'inherit'})

//docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v %CD%\tests\test1-save-two-containers\tmp:/backup dback dback