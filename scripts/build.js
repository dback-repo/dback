const mv = require('fs').renameSync
const del = require('fs').unlinkSync
const cmd = require('child_process').execSync
const r = process.env.REPO


process.env.GOOS = 'linux'
process.env.CGO_ENABLED = '0'

//cmd('docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock dback dback',{stdio: 'inherit'})

//cmd('go build -a -installsuffix cgo -ldflags="-s -w"', {cwd: r+'/go-app/src/dback'})
//cmd('upx --brute dback', {cwd: r+'/go-app/src/dback'})
cmd('go build -a -installsuffix cgo', {cwd: r+'/go-app/src/dback'})
mv(r+'/go-app/src/dback/dback', r+'/docker/dback')
cmd('docker build -t dback .', {cwd: r+'/docker'})
del(r+'/docker/dback')