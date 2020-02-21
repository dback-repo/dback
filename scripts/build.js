const t = new(require('./tools/tools.js'));
const r = process.env.REPO

if (!process.argv[2]){
	process.argv[2] = 'dev'
}

t.cmd('docker build -t dback --target '+process.argv[2]+' .', {cwd: r})