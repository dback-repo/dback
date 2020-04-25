package logic

import (
	"dback/utils/dockerwrapper"
	"log"
)

func Backup(dockerWrapper *dockerwrapper.DockerWrapper, emulate []string, exclude []string) {
	log.Println(`Backup started`)
}
