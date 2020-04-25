package dockerbuilder

import (
	"log"

	"github.com/docker/docker/client"
)

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

func NewDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	check(err, `Cannot connect to docker.
Check:
- dback runs inside docker container, with passed docker socket 
docker run -t --rm -v //var/run/docker.sock:/var/run/docker.sock dback .......
TODO: more details`)

	return cli
}
