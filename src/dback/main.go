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
	spaceTracker := spacetracker.NewSpaceTracker(time.Second)

	switch cliRequest.Command {
	case `backup`:
		logic.Backup(dockerWrapper, dbackOpts, resticWrapper, spaceTracker)
	case `ls`:
		s3Wrapper := s3wrapper.NewS3Wrapper(resticOpts.S3Opts)
		logic.List(s3Wrapper, resticWrapper, dockerWrapper, dbackArgs, spaceTracker)
	case `restore`:
		s3Wrapper := s3wrapper.NewS3Wrapper(resticOpts.S3Opts)
		logic.Restore(s3Wrapper, resticWrapper, dockerWrapper, dbackOpts, dbackArgs, spaceTracker)
	case `inspect`:
		logic.Inspect(dockerWrapper, dbackArgs)
	}
}
