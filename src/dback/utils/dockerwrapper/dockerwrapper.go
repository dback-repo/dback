package dockerwrapper

import (
	"log"

	"github.com/docker/docker/client"
)

type DockerWrapper struct {
	Cli *client.Client
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

func (t *DockerWrapper) Close() {
	check(t.Cli.Close(), `Cannot close connection docker connection`)
}
