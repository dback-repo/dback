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