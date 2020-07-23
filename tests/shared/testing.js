const r = process.env.REPO
const t = new(require(r+'/scripts/tools/tools.js'));
var fs  = require('fs');
const Path = require('path');


module.exports = class Testing {
	constructor()
    {
  //       new(require('../env.js'));
		// this.mv = require('fs').renameSync
		// this.del = require('fs').unlinkSync
		// this.ccmd = require('child_process').execSync
		// this.r = process.env.REPO
		this.createFolder = function(dir) {
			if (!fs.existsSync(dir)){
				fs.mkdirSync(dir);
			}
		}

    }
	// createFolder = function(dir) {
	// 	if (!fs.existsSync(dir)){
	// 		fs.mkdirSync(dir);
	// 	}
	// }

	//one does not simply copy file to the volume!
	//we must mount the volume to temporary container
	moveFileToDockerVolume = function(file, volume){
		t.cmd('docker run --rm -v '+Path.resolve(file)+':/'+Path.basename(file)+' -v '+volume+':/dest busybox cp -r /'+Path.basename(file)+' /dest')
	}

	//move an old tmp folder to tmp-old
	//because of deleting is extrimely slow for multiple files
	clearTmp = function(){
		this.createFolder('tmp')
		this.createFolder('tmp-old')
		fs.renameSync('tmp', 'tmp-old/'+Math.floor(new Date() / 1000))
		this.createFolder('tmp')
	}

	initVolumeWithFile = function(volumeName, filename){
		//try {cmd('docker volume rm ' + volumeName)}catch{}
		t.cmd('docker volume create ' + volumeName)
		this.moveFileToDockerVolume(filename, volumeName)
	}

	checkSub = function(str, substr){
		if (!str.includes(substr)){
			throw('"'+substr+'" is expected but not found in the string:'+"\r\n"+str)
		}
	}

	checkNoSub = function(str, substr){
		if (str.includes(substr)){
			throw('"'+substr+'" is expected but not found in the string:'+"\r\n"+str)
		}
	}

 //    cmd(c, options) {
 //    	try{
 //    		//execSync = require('child_process').execSync
 //    		return this.ccmd(c, options)
 //    	}catch(e){
 //    		console.log('==========================Command failed: '+c)
 //    		console.log('==========================Command output=============================')
 //       		if (e.stdout){
 //    			console.log(e.stdout.toString())
 //    		}
 //    		if (e.stack){
	//     		console.log(e.stack)
	//     	}
 //    		console.log('==========================/Command output=============================')
 //    		throw 'cmd failed'
 //    	}
 //    }

 //  	//env must be stored at ../env.js
	// startCmdDetached(cmd, args) {
	// 	var out = require('fs').openSync('.', 'a')
	// 	var child = require('child_process').spawn(cmd, args, { detached: true, stdio: [ 'ignore', out, out ] })
	// 	child.unref()
	// }

	// checkCmdAvailable(cmd) {
	// 	try {
	// 		this.cmd(cmd)
	// 	}catch{
	// 		console.log('The command "'+cmd+'" returned non-zero code. Check it installed and available')
	// 		console.log('Press any key to exit');
			
	// 		require('fs').readSync(process.stdin.fd, new Buffer(1), 0, 1)
	// 		throw 'check cmd available failed'
	// 	}
	// }
};
