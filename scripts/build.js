const mv = require('fs').renameSync;
const cmd = require('child_process').execSync;

process.env.GOOS = 'linux'

cmd('go build', {cwd: process.env.REPO+'/go-app/src/dback'});
mv(process.env.REPO+'/go-app/src/dback/dback', process.env.REPO+'/docker/dback'); 
