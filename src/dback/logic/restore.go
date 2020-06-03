package logic

import (
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"log"
)

func Restore(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper) {
	log.Println(`Make restore...`)

	s3SavedContainers := s3.GetBackupsContainerList()

	for _, curContainerName := range s3SavedContainers {
		log.Println(curContainerName)
		resticw.ListSnapshots(`/dback-test-1.1/mount-dir/`)
	}
}
