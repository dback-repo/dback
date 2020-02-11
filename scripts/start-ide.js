// process.env['PATH']+=';'+__dirname+'\\..\\node_modules\\liteide\\bin';
t = require('./tools/tools.js')
var tool = new t;
console.log(tool.type)
tool.test()

// var out = require('fs').openSync('.', 'a')
// var child = require('child_process').spawn('cmd', ['/c start cmd'], { detached: true, stdio: [ 'ignore', out, out ] })
// child.unref()

// var out = require('fs').openSync('.', 'a')
// var child = require('child_process').spawn('cmd', ['/c start liteide'], { detached: true, stdio: [ 'ignore', out, out ] })
// child.unref()

//Del \\?\D:\PET\dback2\NUL
