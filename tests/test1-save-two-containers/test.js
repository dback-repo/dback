const mv = require('fs').renameSync
const cmd = require('child_process').execSync
const r = process.env.REPO

try {cmd('docker rm -f dback-test1.1 dback-test1.2')}catch{}
try {cmd('docker volume rm dback-test1.1')}catch{}
cmd('docker volume create dback-test1.1')

cmd('docker run -d --rm --name dback-test1.1 -v %CD%\\data\\mount-dir:/mount docker:19.03.5 tail -f /dev/null')
cmd('docker run -d --rm --name dback-test1.2 -v %CD%\\data\\mount-dir:/mount docker:19.03.5 tail -f /dev/null')
cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock dback dback',{stdio: 'inherit'})

