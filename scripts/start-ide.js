const t = new(require('./tools/tools.js'))

t.checkCmdAvailable('docker ps')

t.startCmdDetached('cmd',['/c start cmd /k echo Hi! Run \'npm run iter\' for build and test locally. Update something in IDE, check it correct (ctrl+b), and run iter again.'])
t.startCmdDetached('liteide', [process.env['REPO']+'/src/dback/main.go'])