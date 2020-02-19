const t = new(require('./tools/tools.js'))

t.checkCmdAvailable('docker ps')
t.cmd('npm i')

t.startCmdDetached('cmd',['/c start cmd /k echo Hi! Run \'npm run iter\' for build and test locally. Update something in IDE, check it correct (ctrl+b), and run iter again.'])
t.startCmdDetached('liteide', [process.env['REPO']+'/go-app/src/dback/main.go'])