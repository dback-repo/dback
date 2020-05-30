module.exports = class SecEnv {
	constructor()
    {
        process.env['DBACK_DOCKER_LOGIN']       = ''
        process.env['DBACK_DOCKER_PASSWORD']    = ''
        process.env['DBACK_DOCKER_REPO']    	= ''
    }
};