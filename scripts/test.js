const t = new(require('./tools/tools.js'))

const r = process.env.REPO

t.cmd('node test.js', {cwd: r+'/tests/test1-save-two-containers',stdio: 'inherit'});