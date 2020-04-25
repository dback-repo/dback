package logic

import (
	"dback/utils/dockerwrapper"
	"log"
)

func Backup(dockerWrapper *dockerwrapper.DockerWrapper, isEmulation dockerwrapper.EmulateFlag,
	excludePatterns []dockerwrapper.ExcludePattern) {
	log.Println(`Backup started`)
	dockerWrapper.GetContanersAndMountsForBackup(excludePatterns)
}
