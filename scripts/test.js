const mv = require('fs').renameSync;
const cmd = require('child_process').execSync;
const r = process.env.REPO

cmd('node test.js', {cwd: r+'/tests/test1-save-two-containers',stdio: 'inherit'});