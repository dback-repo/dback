package logic

import (
	"dback/utils/s3wrapper"
	"log"
)

func Restore(s3 *s3wrapper.S3Wrapper) {
	log.Println(`Make restore...`)
	log.Println(s3.GetBackupsContainerList())
}
