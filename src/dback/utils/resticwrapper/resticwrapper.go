package resticwrapper

import (
	"log"
	"os"
	"os/exec"
)

type ResticWrapper struct {
	s3Endpoint     string
	s3Bucket       string
	accKey         string
	secKey         string
	resticPassword string
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg + "\r\n" + err.Error())
	}
}

func NewResticWrapper(s3Endpoint, s3Bucket, accKey, secKey, resticPassword string) *ResticWrapper {
	return &ResticWrapper{
		s3Endpoint:     s3Endpoint,
		s3Bucket:       s3Bucket,
		accKey:         accKey,
		secKey:         secKey,
		resticPassword: resticPassword,
	}
}

func (t *ResticWrapper) cmd(localFolder, s3Folder, command string, arg ...string) (string, error) {
	commandArg := append([]string{command}, arg...)

	cmd := exec.Command(`/bin/restic`, commandArg...)
	cmd.Dir = localFolder

	cmd.Env = append(os.Environ(),
		`RESTIC_PASSWORD=`+t.resticPassword,
		`RESTIC_REPOSITORY=s3:http://`+t.s3Endpoint+`/`+t.s3Bucket+s3Folder,
		`AWS_ACCESS_KEY_ID=`+t.accKey,
		`AWS_SECRET_ACCESS_KEY=`+t.secKey)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func (t *ResticWrapper) init(localFolder, s3Folder string) error {
	_, err := t.cmd(localFolder, s3Folder, `init`)
	return err
}

func (t *ResticWrapper) backup(localFolder, s3Folder string) (string, error) {
	return t.cmd(localFolder, s3Folder, `backup`, `/`+localFolder)
}

func (t *ResticWrapper) Save(localFolder, s3Folder string) {
	if t.init(localFolder, s3Folder) != nil { //if repo is already created, or other error
		out, err := t.backup(localFolder, s3Folder)
		check(err, out)
	}
}
