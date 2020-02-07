const mv = require('fs').renameSync
const cmd = require('child_process').execSync
const r = process.env.REPO

try {cmd('docker rm -f dback-test1.1 dback-test1.2')}catch{}
try {cmd('docker volume rm dback-test1.1')}catch{}
cmd('docker volume create dback-test1.1')

cmd('docker run -d --name dback-test1.1 -v %CD%\\data\\mount-dir:/mount nginx:1.17.8-alpine')
cmd('docker run -d --name dback-test1.2 -v %CD%\\data\\mount-dir:/mount nginx:1.17.8-alpine')
cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock -v %CD%\\tmp\\:/backup dback dback',{stdio: 'inherit'})