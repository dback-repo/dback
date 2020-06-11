package main

import (
	"dback/logic"
	"dback/utils/cli"
	"dback/utils/dockerbuilder"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"dback/utils/spacetracker"
	"log"
	"time"
)

func main() {
	cliRequest := cli.ParseCLI()
	dbackOpts, resticOpts := cli.VerifyAndCast(cliRequest)

	dockerWrapper := &dockerwrapper.DockerWrapper{Docker: dockerbuilder.NewDockerClient()}
	defer dockerWrapper.Close()

	resticWrapper := resticwrapper.NewResticWrapper(resticOpts)
	s3Wrapper := s3wrapper.NewS3Wrapper(resticOpts.S3Opts)

	spaceTracker := spacetracker.NewSpaceTracker(time.Second)

	switch cliRequest.Command {
	case `backup`:
		logic.Backup(dockerWrapper, dbackOpts, resticWrapper)
		log.Println(`Minimal disk space: `, spaceTracker.MinSpaceBytes)
		log.Println(`Used space: `, spaceTracker.StartSpace-spaceTracker.MinSpaceBytes)
	case `restore`:
		logic.Restore(s3Wrapper, resticWrapper, dockerWrapper, dbackOpts)
		log.Println(`Minimal disk space: `, spaceTracker.MinSpaceBytes)
		log.Println(`Used space: `, spaceTracker.StartSpace-spaceTracker.MinSpaceBytes)
	case `list`:
	}
}
