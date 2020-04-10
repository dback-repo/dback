package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	var s3Endpoint, s3Bucket, accKey, secKey string

	if runtime.GOOS == "linux" {
		check(os.MkdirAll(`/tmp`, 664))
	}

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
				fmt.Println(`Usage:  dback backup [OPTIONS]

Make snapshot of mounts matched all the points:
default points:
- HostConfig.AutoRemove:      false
- HostConfig.RestartPolicy:   always
- Status.State:               running
- Status.Running:             true

Options:
  --exclude-mount string            Exclude volume pattern
    mounts are named as: [ContainerName]/[PathInContainer]
    For example, mount in "mysql" container: mysql/var/mysql/data
    Pattern is regular expression. For example, "^/(drone.*|dback-test-1.5.*)$"
    ignore all mounts starts with "/drone", or "/dback-test-1.5"`)
				return
			case `--exclude-mount`:
				if len(args) < 3 {
					fmt.Println(`Exclude parameter is defined, but pattern is not provided`)
					return
				}
				excludePattern = args[2]
				s3Endpoint = args[3]
				s3Bucket = args[4]
				accKey = args[5]
				secKey = args[6]

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
			go backupContainer(curContainer, &wg, excludePattern, s3Endpoint, s3Bucket, accKey, secKey)
		}

		wg.Wait()
		log.Println(`Backup has finished for the mounts above`)
	case `restore`:
		log.Println(`Restore started`)
		log.Println(args)
		s3Endpoint = args[1]
		s3Bucket = args[2]
		accKey = args[3]
		secKey = args[4]

		restoreContainers(getBackupsContainerList(s3Endpoint, s3Bucket, accKey, secKey), s3Endpoint, s3Bucket, accKey, secKey)
		log.Println(`Restore has finished for the mounts above`)
	case `help`:
		fmt.Println(`Usage:  dback [OPTIONS] COMMAND

A tool for docker mounts bulk backup and restore

Options:
  --folder string      Not implemented yet. Location of client config files (default "dback-snapshots")
Commands:
  backup               Make snapshot of mounts
  restore              Restore snapshots to exist mounts

Run 'dback COMMAND --help' for more information on a command`)
	default:
		fmt.Println(`Unknown command. Type "dback help" for see manual`)
	}
}
