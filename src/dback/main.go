package main

import (
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
		backup(dockerClient, f[`emulate`], f[`x`])
	case `restore`:
		restore(dockerClient, f[`emulate`])
	case ``: //no command provided. Parse CLI is already printed an advice
		os.Exit(1)
	default:
		log.Fatalln(`Unrecognized command ` + cliRequest.Command +
			`. Run with --help, for list of available commands`)
	}
}
