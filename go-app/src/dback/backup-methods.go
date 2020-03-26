package main

import (
	"context"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
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

func unpackTarToMyself(c types.Container, myContainer types.Container, m types.MountPoint) {
	//check(os.MkdirAll(`/tmp`, 664))

	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	tar, err := os.Open(`dback-snapshots/` + c.Names[0] + m.Destination + `/tar.tar`)
	check(err)
	defer tar.Close()

	lastSlashIdx := strings.LastIndex(m.Destination, `/`)
	destParent := m.Destination[:lastSlashIdx] //      "/var/www/lynx" -> "/var/www"        "/opt" -> "/"
	if destParent == `` {
		destParent = `/`
	}

	destParent = `/dback-snapshots` + c.Names[0] + m.Destination

	check(cli.CopyToContainer(context.Background(), myContainer.ID, destParent, tar, types.CopyToContainerOptions{true, false}))
}

func backupMount(cli *client.Client, c types.Container, m types.MountPoint, wg *sync.WaitGroup) {
	defer wg.Done()
	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	check(os.MkdirAll(`dback-snapshots/`+c.Names[0]+m.Destination, 0664))

	reader, _, err := cli.CopyFromContainer(context.Background(), c.ID, m.Destination)
	check(err)

	outFile, err := os.Create(`dback-snapshots/` + c.Names[0] + m.Destination + `/tar.tar`)
	check(err)
	defer outFile.Close()
	_, err = io.Copy(outFile, reader)
	check(err)

	log.Println(`make backup: ` + c.Names[0] + m.Destination)

	myselfContainerID, err := os.Hostname()
	check(err)

	unpackTarToMyself(c, *getContainerByNameOrId(cli, myselfContainerID), m)
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
				if !inspect.HostConfig.RestartPolicy.IsNone() {
					mounts := []types.MountPoint{}
					for _, curMount := range c.Mounts {
						if excludePattern == `` {
							mounts = append(mounts, curMount)
						} else {
							if !r.MatchString(c.Names[0] + curMount.Destination) {
								mounts = append(mounts, curMount)
							} else {
								log.Println(`exclude: ` + c.Names[0] + curMount.Destination + `      Reason: --exclude-mount parameter`)
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
						go backupMount(cli, c, curMount, &wgMount)
					}
					wgMount.Wait()

					check(cli.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{}))
				} else {
					log.Println(`exclude: ` + c.Names[0] + `      Reason: container restart policy==none`)
				}
			} else {
				log.Println(`exclude: ` + c.Names[0] + `      Reason: temporary container (--rm)`)
			}
		}
	}
}
