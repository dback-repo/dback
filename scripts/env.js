module.exports = class Env {
	constructor()
    {
		if (! "ENVISSET" in process.env) {
			return
		}

    	var sep=':'
		if (process.platform === "win32"){
			sep=';'
		}

    	process.env['REPO'] = 					require('path').resolve(__dirname+'/..')

        process.env['GOPATH'] = 				process.env['REPO']+'/go-app'
        process.env['GOROOT'] = 				process.env['REPO']+'/node_modules/go-win'

        process.env['PATH'] += 					sep+process.env['GOPATH']+'/bin'
        process.env['PATH'] += 					sep+process.env['GOROOT']+'/bin'
        process.env['PATH'] += 					sep+process.env['REPO']+'/node_modules/.bin'
        process.env['PATH'] += 					sep+process.env['REPO']+'/node_modules/upx-win'
        process.env['PATH'] += 					sep+process.env['REPO']+'/node_modules/liteide/bin'

        process.env['DOCKER_API_VERSION'] = 	'1.37'
        process.env['ENVISSET'] = 				'TRUE'
    }
};