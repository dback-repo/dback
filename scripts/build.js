const t = new(require('./tools/tools.js'));
const r = process.env.REPO

switch (process.argv[2]) {
case 'dev':
	t.cmd('docker build -t dback -f dockerfile_dev .', {cwd: r})
	break;
case 'prod':
	t.cmd('docker build -t dback .', {cwd: r})
	break;
default:
	console.log(process.argv)
	throw ('Unknown build profile')
}



//process.env.GOOS = 'linux'
//process.env.CGO_ENABLED = '0'

// cmd('go build -a -installsuffix cgo -ldflags="-s -w"', {cwd: r+'/go-app/src/dback'})
// cmd('upx --brute dback', {cwd: r+'/go-app/src/dback'})

//t.cmd('go build -a -installsuffix cgo', {cwd: r+'/go-app/src/dback'})
//t.mv(r+'/go-app/src/dback/dback', r+'/docker/dback')
//t.cmd('docker build -t dback .', {cwd: r+'/docker'})
//t.del(r+'/docker/dback')