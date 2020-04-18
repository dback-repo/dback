module.exports = class Env {
	constructor()
    {
		if (! "ENVISSET" in process.env) {
			return
		}

        process.env['DBACK_VER'] =              '0.0'

    	var sep=':' 
		if (process.platform === "win32"){ // override separator in PATH variable for windows
			sep=';'
		}

    	process.env['REPO'] = 					require('path').resolve(__dirname+'/..')

        process.env['GOPATH'] = 				process.env['REPO']
        process.env['GOROOT'] = 				process.env['REPO']+'/node_modules/go-win'
		if (process.platform === "linux"){
        	process.env['GOROOT'] = 			process.env['REPO']+'/node_modules/go-linux'
		}

        process.env['PATH'] += 					sep+process.env['GOPATH']+'/bin'
        process.env['PATH'] += 					sep+process.env['GOROOT']+'/bin'
        process.env['PATH'] += 					sep+process.env['REPO']+'/node_modules/.bin'
        process.env['PATH'] +=                  sep+process.env['REPO']+'/node_modules/golangci-lint'
        process.env['PATH'] += 					sep+process.env['REPO']+'/node_modules/upx-win'
        process.env['PATH'] += 					sep+process.env['REPO']+'/node_modules/upx-linux'        
        process.env['PATH'] += 					sep+process.env['REPO']+'/node_modules/liteide-win/bin'

        process.env['DOCKER_API_VERSION'] = 	'1.37'
        process.env['ENVISSET'] = 				'TRUE'

        const fs = require('fs')
        if (!process.env['CI']) {                            //load secrets from /secrets folder, if we are not in CI
            if (!fs.existsSync(process.env['REPO']+'/secrets')){    //create secrets file from draft, if not exist
                fs.mkdirSync(process.env['REPO']+'/secrets');
                if (!fs.existsSync(process.env['REPO']+'/secrets/env.js')){
                    const mv = require('fs').copyFileSync
                    mv(process.env['REPO']+'/scripts/tools/secret-env-draft.js',process.env['REPO']+'/secrets/env.js')
                }
            }
            new(require(process.env['REPO']+'/secrets/env.js'));
        }
    }
};