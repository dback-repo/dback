const t = new(require('./tools/tools.js'));
const r = process.env.REPO

if(!process.env.DRONE_BUILD_NUMBER) {
	throw 'DRONE_BUILD_NUMBER env variable is not set. Set manual if required. Warning: no any protection for exist build numbers'
}

t.cmd('docker login -u '+process.env.DBACK_DOCKER_LOGIN+' -p '+process.env.DBACK_DOCKER_PASSWORD, {cwd: r})
t.cmd('docker tag dback:latest '+process.env.DBACK_DOCKER_LOGIN+'/'+process.env.DBACK_DOCKER_REPO+':latest', {cwd: r})
t.cmd('docker tag dback:latest '+process.env.DBACK_DOCKER_LOGIN+'/'+process.env.DBACK_DOCKER_REPO+':'+process.env.DBACK_VER+'.'+process.env.DRONE_BUILD_NUMBER, {cwd: r})
t.cmd('docker push '+process.env.DBACK_DOCKER_LOGIN+'/'+process.env.DBACK_DOCKER_REPO+':latest', {cwd: r})
t.cmd('docker push '+process.env.DBACK_DOCKER_LOGIN+'/'+process.env.DBACK_DOCKER_REPO+':'+process.env.DBACK_VER+'.'+process.env.DRONE_BUILD_NUMBER, {cwd: r})