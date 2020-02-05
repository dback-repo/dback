const mv = require('fs').renameSync;
const cmd = require('child_process').execSync;
const r = process.env.REPO

//cmd('go build', {cwd: r+'/go-app/src/dback'});
//mv(r+'/go-app/src/dback/dback', r+'/docker/dback');

cmd('node test.js', {cwd: r+'/tests/test1-save-two-containers'});