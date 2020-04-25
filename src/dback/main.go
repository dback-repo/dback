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
	cliRequest := cli.ParseCLI() //will be interrupted with printing adivce here, on parsing error

	dockerWrapper := &dockerwrapper.DockerWrapper{Cli: dockerbuilder.NewDockerClient()}

	f := cliRequest.Flags

	switch cliRequest.Command {
	case `backup`:
		logic.Backup(dockerWrapper, f[`emulate`], f[`x`])
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
