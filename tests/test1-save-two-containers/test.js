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

//move an old tmp folder to tmp-old
//because of deleting is extrimely slow for multiple files
createFolder('tmp-old')
fs.renameSync('tmp', 'tmp-old/'+Math.floor(new Date() / 1000))
createFolder('tmp')

try {cmd('docker rm -f dback-test-1.1 dback-test-1.2')}catch{}
try {cmd('docker volume rm dback-test-1.1')}catch{}
cmd('docker volume create dback-test-1.1')

cmd('docker run -d --name dback-test-1.1 -v %CD%\\data\\mount-dir:/mount nginx:1.17.8-alpine')
cmd('docker run -d --name dback-test-1.2 -v %CD%\\data\\mount-dir:/mount nginx:1.17.8-alpine')
cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v %CD%\\tmp:/backup dback dback',{stdio: 'inherit'})

//docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v %CD%\tests\test1-save-two-containers\tmp:/backup dback dback