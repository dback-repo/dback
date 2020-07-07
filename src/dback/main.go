package main

import (
	"dback/logic"
	"dback/utils/cli"
	"dback/utils/dockerbuilder"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"dback/utils/spacetracker"
	"time"
)

func main() {
	cliRequest := cli.ParseCLI()
	dbackOpts, dbackArgs, resticOpts := cli.VerifyAndCast(cliRequest)

	dockerWrapper := &dockerwrapper.DockerWrapper{Docker: dockerbuilder.NewDockerClient()}
	defer dockerWrapper.Close()

	resticWrapper := resticwrapper.NewResticWrapper(resticOpts)
	s3Wrapper := s3wrapper.NewS3Wrapper(resticOpts.S3Opts)

	spaceTracker := spacetracker.NewSpaceTracker(time.Second)

	switch cliRequest.Command {
	case `backup`:
		logic.Backup(dockerWrapper, dbackOpts, resticWrapper, spaceTracker)
	case `ls`:
		logic.List(s3Wrapper, resticWrapper, dockerWrapper, spaceTracker)
	case `restore`:
		logic.Restore(s3Wrapper, resticWrapper, dockerWrapper, dbackOpts, dbackArgs, spaceTracker)
	}
}
