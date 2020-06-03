package resticwrapper

import (
	"dback/utils/s3wrapper"
	"log"
	"os"
	"os/exec"
)

type CreationOpts struct {
	S3Opts     s3wrapper.CreationOpts
	ResticPass string
}

type ResticWrapper struct {
	s3Endpoint string
	s3Bucket   string
	accKey     string
	secKey     string
	resticPass string
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg + "\r\n" + err.Error())
	}
}

func NewResticWrapper(opts CreationOpts) *ResticWrapper {
	return &ResticWrapper{
		s3Endpoint: opts.S3Opts.S3Endpoint,
		s3Bucket:   opts.S3Opts.S3Bucket,
		accKey:     opts.S3Opts.AccKey,
		secKey:     opts.S3Opts.SecKey,
		resticPass: opts.ResticPass,
	}
}

func (t *ResticWrapper) cmd(localFolder, s3Folder, command string, arg ...string) (string, error) {
	commandArg := append([]string{command}, arg...)

	cmd := exec.Command(`/bin/restic`, commandArg...)
	cmd.Dir = localFolder

	cmd.Env = append(os.Environ(),
		`RESTIC_PASSWORD=`+t.resticPass,
		`RESTIC_REPOSITORY=s3:`+t.s3Endpoint+`/`+t.s3Bucket+s3Folder,
		`AWS_ACCESS_KEY_ID=`+t.accKey,
		`AWS_SECRET_ACCESS_KEY=`+t.secKey)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func (t *ResticWrapper) init(localFolder, s3Folder, tag string) error {
	_, err := t.cmd(localFolder, s3Folder, `init`)
	return err
}

func (t *ResticWrapper) backup(localFolder, s3Folder, tag string) (string, error) {
	return t.cmd(localFolder, s3Folder, `backup`, `--tag`, tag, `/`+localFolder)
}

func (t *ResticWrapper) Save(localFolder, s3Folder, tag string) {
	if t.init(localFolder, s3Folder, tag) != nil { //if repo is already created, or other error
		out, err := t.backup(localFolder, s3Folder, tag)
		check(err, out)
	}
}
