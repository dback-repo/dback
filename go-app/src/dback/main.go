package main

import (
	"context"
	"log"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	log.Println(`Backup started`)
	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	check(err)

	var wg sync.WaitGroup
	wg.Add(len(containers))

	for _, curContainer := range containers {
		go backupContainer(curContainer, &wg)
	}

	wg.Wait()
	log.Println(`Backup success`)
}
