package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		args = append(args, `help`)
	}

	switch args[0] {
	case `backup`:
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
		log.Println(`Backup has finished for the mounts above`)
	case `restore`:
		log.Println(`Restore started`)
		restoreContainers(getBackupsContainerList())
		log.Println(`Restore has finished for the mounts above`)
	case `help`:
		fmt.Println("Here is no manual yet...   :(")
	default:
		fmt.Println("Unknown command")
	}
}
