module.exports = class Tools {
	constructor()
    {
        new(require('../env.js'));
    }

  	//env must be stored at ../env.js
	startCmdDetached(cmd) {
		var out = require('fs').openSync('.', 'a')
		var child = require('child_process').spawn('cmd', ['/c start '+cmd], { detached: true, stdio: [ 'ignore', out, out ] })
		child.unref()
	}
};
