
//process.env['PATH']+=';'+__dirname+'\\..\\node_modules\\liteide\\bin';

// var out = require('fs').openSync('.', 'a')
// var child = require('child_process').spawn('cmd', ['/c start cmd'], { detached: true, stdio: [ 'ignore', out, out ] })
// child.unref()

//Del \\?\D:\PET\dback2\NUL

module.exports = class Tools {
	constructor()
    {
        this.type = "Tools";
    }

  test() {
  	console.log(this.type);
    //return this.width ** 2;
  }
};
