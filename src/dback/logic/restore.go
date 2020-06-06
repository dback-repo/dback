package logic

import (
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"log"
)

func getS3MountsForRestore(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, restoreParams []string) []s3wrapper.S3Mount {
	s3Mounts := s3.GetMountsForRestore(restoreParams, resticw, dockerw)
	res := s3Mounts

	return res
}

func Restore(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, restoreParams []string) {
	s3MountsForRestore := getS3MountsForRestore(s3, resticw, dockerw, restoreParams)

	log.Println(s3MountsForRestore)
}
