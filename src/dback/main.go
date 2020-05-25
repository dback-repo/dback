package main

import (
	"dback/logic"
	"dback/utils/cli"
	"dback/utils/dockerbuilder"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"log"
	"os"
)

func main() {
	cliRequest := cli.ParseCLI()
	isEmulation, excludePatterns := cli.VerifyAndCast(cliRequest)

	dockerWrapper := &dockerwrapper.DockerWrapper{Docker: dockerbuilder.NewDockerClient()}
	defer dockerWrapper.Close()

	resticWrapper := resticwrapper.NewResticWrapper(``, ``, ``, ``, ``)

	switch cliRequest.Command {
	case `backup`:
		logic.Backup(dockerWrapper, isEmulation, excludePatterns, 1, resticWrapper)
	case `restore`:
		logic.Restore()
	case ``:
		//no command provided. Parse CLI is already printed an advice
		os.Exit(1)
	default:
		log.Fatalln(`Unrecognized command ` + cliRequest.Command +
			`. Run with --help, for list of available commands`)
	}
}
