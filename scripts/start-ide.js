const t = new(require('./tools/tools.js'));
t.startCmdDetached('cmd',['/c start cmd'])
t.startCmdDetached('liteide', [process.env['REPO']+'/go-app/src/dback/main.go'])