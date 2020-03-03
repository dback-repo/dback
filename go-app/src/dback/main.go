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
		excludePattern := ``
		if len(args) > 1 {
			switch args[1] {
			case `--help`:
				fmt.Println(`backup help`)
				return
			case `--exclude-mount`:
				if len(args) < 3 {
					fmt.Println(`Exclude parameter is defined, but pattern is not provided`)
					return
				}
				excludePattern = args[2]
			default:
				fmt.Println("Unknown parameter")
				return
			}
		}

		log.Println(`Backup started`)
		cli, err := client.NewEnvClient()
		check(err)
		defer cli.Close()

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		check(err)

		var wg sync.WaitGroup
		wg.Add(len(containers))

		for _, curContainer := range containers {
			go backupContainer(curContainer, &wg, excludePattern)
		}

		wg.Wait()
		log.Println(`Backup has finished for the mounts above`)
	case `restore`:
		log.Println(`Restore started`)
		restoreContainers(getBackupsContainerList())
		log.Println(`Restore has finished for the mounts above`)
	case `help`:
		fmt.Println("Here is no manual yet.....   :(")
	default:
		fmt.Println("Unknown command")
	}
}
