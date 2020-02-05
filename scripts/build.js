const mv = require('fs').renameSync;
const cmd = require('child_process').execSync;
const r = process.env.REPO


process.env.GOOS = 'linux'

//cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock dback dback',{stdio: 'inherit'})

cmd('go build', {cwd: r+'/go-app/src/dback'});
mv(r+'/go-app/src/dback/dback', r+'/docker/dback');
cmd('docker build -t dback .', {cwd: r+'/docker'});