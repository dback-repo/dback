package logic

import (
	"dback/utils/cli"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"dback/utils/spacetracker"
	"log"
)

func restoreContainer(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, dbackOpts cli.DbackOpts, spacetracker *spacetracker.SpaceTracker) {
	log.Println(`Restore container...`)
	log.Println(`Restore finished for the mounts above`)
}
