const t = new(require('./tools/tools.js'));
const r = process.env.REPO

t.cmd('golangci-lint run --enable-all', {cwd: r+'/src/dback'})