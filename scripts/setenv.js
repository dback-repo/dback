// const mv = require('fs').renameSync;
// const cmd = require('child_process').execSync;

// const r = process.env.REPO

// console.log('setenv');
// console.log(process.env.REPO);
// //sdf
// cmd('cmd.exe /c ping ya.ru')
// process.env.REPO = 'testValue'
// console.log(process.env.REPO);

//cmd('node test.js', {cwd: r+'/tests/test1-save-two-containers',stdio: 'inherit'});

// On Windows Only...
// const { spawn } = require('child_process')
// const bat = spawn('cmd.exe', ['/c', 'start cmd'])

// bat.stdout.on('data', (data) => {
//   console.log(data.toString());
// });

// bat.stderr.on('data', (data) => {
//   console.error(data.toString());
// });

// bat.on('exit', (code) => {
//   console.log(`Child exited with code ${code}`);
// });

// const spawn = require('child_process').spawn;
// console.log(process.argv[0])
// const child = spawn('cmd', ['/c start cmd'], {
//   detached: true,
//   stdio: ['ignore']
// });

// child.unref();






// const child_process = require('child_process');
// var p = child_process.fork('cmd.exe /c start cmd')
// p.disconnect(); 
// p.unref();








var fs = require('fs')
var out = fs.openSync('.', 'a')
var cp = require('child_process')
var child = cp.spawn('cmd', ['/c start cmd'], { detached: true, stdio: [ 'ignore', out, out ] })
child.unref()