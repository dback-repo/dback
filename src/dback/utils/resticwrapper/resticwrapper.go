package resticwrapper

import (
	"dback/utils/s3opts"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CreationOpts struct {
	S3Opts     s3opts.CreationOpts
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
	err := t.init(localFolder, s3Folder, tag)
	if err == nil { //ignore error (it just because repo already exist), and prevent linter alerting.
		fmt.Print(``)
	}

	out, err := t.backup(localFolder, s3Folder, tag)

	check(err, out)
}

type Snapshot struct {
	ID  string
	Tag string
}

const snapshotDataLinesOffset = 2

// created new cache in /.cache/restic
// ID        Time                 Host          Tags                 Paths
// ---------------------------------------------------------------------------------------------------------------------
// 2d54ed46  2020-06-06 09:54:31  63a47ed66844  06.06.2020.09-54-27  /tmp/dback-data/mount-data/dback-test-1.1/mount-dir
// 322fe600  2020-06-06 09:54:36  9d54964fad7b  06.06.2020.09-54-35  /tmp/dback-data/mount-data/dback-test-1.1/mount-dir
// ---------------------------------------------------------------------------------------------------------------------
func parseSnapshots(out string) []Snapshot {
	res := []Snapshot{}
	lines := strings.Split(out, "\n")

	headerLine := -1

	for idx, curLine := range lines {
		if strings.HasPrefix(curLine, `ID`) {
			headerLine = idx
			break
		}
	}

	if headerLine == -1 {
		return res
	}

	tagsOffset := strings.Index(lines[headerLine], `Tags`)

	for i := headerLine + snapshotDataLinesOffset; i <= len(lines); i++ {
		if string(lines[i][0]) == `-` {
			break
		}

		tag := lines[i][tagsOffset:]        //06.06.2020.09-54-35  /tmp/d...
		tag = tag[:strings.Index(tag, ` `)] //06.06.2020.09-54-35
		snapshot := Snapshot{}
		snapshot.ID = lines[i][:strings.Index(lines[i], ` `)]
		snapshot.Tag = tag
		res = append(res, snapshot)
	}

	return res
}

func (t *ResticWrapper) ListSnapshots(s3Folder string) []Snapshot {
	out, err := t.cmd(`.`, s3Folder, `snapshots`)
	check(err, `cannot execute "restic snapshots"`)

	return parseSnapshots(out)
}

func (t *ResticWrapper) restore(localFolder, s3Folder, snapshotID string) (string, error) {
	return t.cmd(localFolder, s3Folder, `restore`, snapshotID, `--target`, localFolder)
}

func (t *ResticWrapper) Load(localFolder, s3Folder, snapshotID string) {
	out, err := t.restore(localFolder, s3Folder, snapshotID)
	check(err, out)
}
