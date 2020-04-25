package main

import (
	"dback/logic"
	"dback/utils/cli"
	"dback/utils/dockerbuilder"
	"dback/utils/dockerwrapper"
	"log"
	"os"
)

func main() {
	cliRequest := cli.ParseCLI()
	isEmulation, excludePatterns := verifyCliReq(cliRequest)

	dockerWrapper := &dockerwrapper.DockerWrapper{Docker: dockerbuilder.NewDockerClient()}
	defer dockerWrapper.Close()

	switch cliRequest.Command {
	case `backup`:
		logic.Backup(dockerWrapper, isEmulation, excludePatterns)
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
