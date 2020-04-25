package dockerwrapper

import (
	"log"

	"github.com/docker/docker/client"
)

type DockerWrapper struct {
	Docker *client.Client
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

func (t *DockerWrapper) Close() {
	check(t.Docker.Close(), `Cannot close connection docker connection`)
}

func (t *DockerWrapper) GetContanersAndMountsForBackup(excludePatterns []ExcludePattern) {
	//check(t.docker.Close(), `Cannot close connection docker connection`)
}
