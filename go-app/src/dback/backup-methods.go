package main

import (
	"context"
	"io"
	"log"
	"os"
	"regexp"

	//	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func backupMount(c types.Container, m types.MountPoint, wg *sync.WaitGroup) {
	defer wg.Done()
	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	check(os.MkdirAll(`/backup/`+c.Names[0]+m.Destination, 0664))

	reader, _, err := cli.CopyFromContainer(context.Background(), c.ID, m.Destination)
	check(err)

	outFile, err := os.Create(`/backup/` + c.Names[0] + m.Destination + `/tar.tar`)
	check(err)
	defer outFile.Close()
	_, err = io.Copy(outFile, reader)

	// reader, _, err := cli.CopyFromContainer(context.Background(), c.ID, m.Destination)
	// check(err)

	// lastSlashIdx := strings.LastIndex(m.Destination, `/`)
	// destParent := m.Destination[:lastSlashIdx] // /var/www/lynx -> /var/www
	// if destParent == `` {
	// 	destParent = `/`
	// }

	// check(Untar(reader, `/backup/`+c.Names[0]+destParent))
	log.Println(`make backup: ` + c.Names[0] + m.Destination)
}

func backupContainer(c types.Container, wg *sync.WaitGroup, excludePattern string) {
	defer wg.Done()

	r, err := regexp.Compile(excludePattern)
	check(err)

	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	if c.State == `running` {
		if len(c.Mounts) > 0 {
			inspect, err := cli.ContainerInspect(context.Background(), c.ID)
			check(err)

			if inspect.HostConfig.AutoRemove == false {
				mounts := []types.MountPoint{}
				for _, curMount := range c.Mounts {
					if excludePattern == `` {
						mounts = append(mounts, curMount)
					} else {
						if !r.MatchString(c.Names[0] + curMount.Destination) {
							mounts = append(mounts, curMount)
						} else {
							log.Println(`exclude: ` + c.Names[0] + curMount.Destination)
						}
					}
				}

				timeout := time.Minute

				if len(mounts) == 0 {
					return
				}
				check(cli.ContainerStop(context.Background(), c.ID, &timeout))

				var wgMount sync.WaitGroup
				wgMount.Add(len(mounts))
				for _, curMount := range mounts {
					go backupMount(c, curMount, &wgMount)
				}
				wgMount.Wait()

				check(cli.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{}))
			}
		}
	}
}
