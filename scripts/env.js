module.exports = class Env {
	constructor()
    {
    	var sep=':'
		if (process.platform === "win32"){
			sep=';'
		}

    	process.env['REPO'] = 		require('path').resolve(__dirname+'/..')
        process.env['PATH'] += 		sep+process.env['REPO']+'/node_modules/liteide/bin'
    }
};