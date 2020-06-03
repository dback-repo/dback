package main

import (
	"dback/logic"
	"dback/utils/cli"
	"dback/utils/dockerbuilder"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
)

func main() {
	cliRequest := cli.ParseCLI()
	dbackOpts, resticOpts := cli.VerifyAndCast(cliRequest)

	dockerWrapper := &dockerwrapper.DockerWrapper{Docker: dockerbuilder.NewDockerClient()}
	defer dockerWrapper.Close()

	resticWrapper := resticwrapper.NewResticWrapper(resticOpts)
	s3Wrapper := s3wrapper.NewS3Wrapper(resticOpts.S3Opts)

	switch cliRequest.Command {
	case `backup`:
		logic.Backup(dockerWrapper, dbackOpts, resticWrapper)
	case `restore`:
		logic.Restore(s3Wrapper)
		// s3Wrapper := NewS3Wrapper
		// logic.List(NewS3Wrapper(resticOpts.S3Opts))
	case `list`:
	}
}
