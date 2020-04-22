package main

import (
	"dback/logic"
	"dback/utils/cli"
	"dback/utils/docker"
	"log"
	"os"
)

func main() {
	cliRequest := cli.ParseCLI() //will be interrupted with printing adivce here, on parsing error

	dockerClient := docker.MustNewDockerClient()
	//dockerClient.CheckWeAreInDocker()

	f := cliRequest.Flags

	switch cliRequest.Command {
	case `backup`:
		logic.Backup(dockerClient, f[`emulate`], f[`x`])
	case `restore`:
		logic.Restore(dockerClient, f[`emulate`])
	case ``: //no command provided. Parse CLI is already printed an advice
		os.Exit(1)
	default:
		log.Fatalln(`Unrecognized command ` + cliRequest.Command +
			`. Run with --help, for list of available commands`)
	}
}
