package resticwrapper

type ResticWrapper struct {
	s3Endpoint     string
	s3Bucket       string
	accKey         string
	secKey         string
	resticPassword string
}

// func check(err error, msg string) {
// 	if err != nil {
// 		log.Fatalln(msg + "\r\n" + err.Error())
// 	}
// }

func NewResticWrapper(s3Endpoint, s3Bucket, accKey, secKey, resticPassword string) *ResticWrapper {
	return &ResticWrapper{
		s3Endpoint:     s3Endpoint,
		s3Bucket:       s3Bucket,
		accKey:         accKey,
		secKey:         secKey,
		resticPassword: resticPassword,
	}
}

func (t *ResticWrapper) Init() {
}

func (t *ResticWrapper) Backup() {
}
