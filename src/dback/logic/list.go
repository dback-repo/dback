package logic

import (
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"dback/utils/spacetracker"
	"log"
)

func List(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, spacetracker *spacetracker.SpaceTracker) {
	s3Mounts := getS3MountsForRestore(s3, resticw, dockerw, ``)

	log.Println(s3Mounts)

}
