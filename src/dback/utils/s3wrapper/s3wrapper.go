package s3wrapper

import (
	"log"

	"github.com/minio/minio-go"
)

type CreationOpts struct {
	S3Endpoint string
	S3Bucket   string
	AccKey     string
	SecKey     string
}

type S3Wrapper struct {
	s3Endpoint string
	s3Bucket   string
	accKey     string
	secKey     string
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg + "\r\n" + err.Error())
	}
}

func NewS3Wrapper(opts CreationOpts) *S3Wrapper {
	return &S3Wrapper{
		s3Endpoint: opts.S3Endpoint,
		s3Bucket:   opts.S3Bucket,
		accKey:     opts.AccKey,
		secKey:     opts.SecKey,
	}
}

//list of saved containers == list of folders in /dback-snapshots
func (t *S3Wrapper) GetBackupsContainerList() []string {
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(t.s3Endpoint, t.accKey, t.secKey, useSSL)
	check(err, `cannot start minio client`)

	doneCh := make(chan struct{})
	defer close(doneCh)

	containers := []string{}

	for object := range minioClient.ListObjects(t.s3Bucket, ``, false, doneCh) {
		check(object.Err, `cannot list object at bucket `+t.s3Bucket)
		containers = append(containers, object.Key[:len(object.Key)-1])
	}

	return containers
}
