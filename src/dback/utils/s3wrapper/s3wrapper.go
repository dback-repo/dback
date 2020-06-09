package s3wrapper

import (
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3opts"
	"log"
	"strings"

	"github.com/minio/minio-go"
)

type S3Wrapper struct {
	minio    *minio.Client
	s3Bucket string
}

type S3Mount struct {
	ContainerName      string
	Dest               string
	Snapshots          []resticwrapper.Snapshot
	SelectedSnapshotID string
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg + "\r\n" + err.Error())
	}
}

func NewS3Wrapper(opts s3opts.CreationOpts) *S3Wrapper {
	useSSL := false
	minioClient, err := minio.New(cutProtocol(opts.S3Endpoint), opts.AccKey, opts.SecKey, useSSL)
	check(err, `cannot start minio client`)

	return &S3Wrapper{minioClient, opts.S3Bucket}
}

//http://host.com -> host.com
func cutProtocol(s3Endpoint string) string {
	s3Endpoint = strings.TrimPrefix(s3Endpoint, `http://`)
	s3Endpoint = strings.TrimPrefix(s3Endpoint, `https://`)

	return s3Endpoint
}

func (t *S3Wrapper) GetMounts(
	resticWrapper *resticwrapper.ResticWrapper, dockerw *dockerwrapper.DockerWrapper) []S3Mount {
	doneCh := make(chan struct{})
	defer close(doneCh)

	resticFolders := []string{}

	for object := range t.minio.ListObjects(t.s3Bucket, ``, true, doneCh) {
		check(object.Err, `cannot list object at bucket `+t.s3Bucket)

		if strings.Contains(object.Key, `/config`) {
			resticFolders = append(resticFolders, object.Key)
		}
	}

	res := []S3Mount{}

	//get restic repos
	for _, curResticFolder := range resticFolders {
		mount := S3Mount{}
		mount.ContainerName = `/` + curResticFolder[:strings.Index(curResticFolder, `/`)]
		mount.Dest = curResticFolder[strings.Index(curResticFolder, `/`):]
		mount.Dest = strings.TrimSuffix(mount.Dest, `/config`)
		res = append(res, mount)
	}

	//get snapshots for each restic repo
	for mountIdx, curMount := range res {
		res[mountIdx].Snapshots = resticWrapper.ListSnapshots(curMount.ContainerName + curMount.Dest)
	}

	return res
}

func GetContainerIDsOfMounts(mounts []S3Mount, dockerW *dockerwrapper.DockerWrapper) []string {
	res := []string{}

	containersMap := make(map[string]string)

	for _, curMount := range mounts {
		containersMap[curMount.ContainerName] = curMount.ContainerName
	}

	for _, curContainerName := range containersMap {
		res = append(res, dockerW.GetContainerIDByName(curContainerName))
	}

	return res
}
