package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
)

//list of saved containers == list of folders in /backup
func getBackupsContainerList() []string {
	var res []string
	files, err := ioutil.ReadDir(`/backup`)
	check(err)
	for _, curFile := range files {
		if curFile.IsDir() {
			res = append(res, curFile.Name())
		}
	}
	return res
}

func restoreContainers(containers []string) {
	var wg sync.WaitGroup
	wg.Add(len(containers))

	for _, curContainer := range containers {
		go restoreContainer(curContainer, &wg)
	}

	wg.Wait()
}

func restoreMount(c types.Container, m types.MountPoint, wg *sync.WaitGroup) {
	defer wg.Done()
	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	tar, err := os.Open(`/backup/` + c.Names[0] + m.Destination + `/tar.tar`)
	check(err)

	lastSlashIdx := strings.LastIndex(m.Destination, `/`)
	destParent := m.Destination[:lastSlashIdx] //      "/var/www/lynx" -> "/var/www"        "/opt" -> "/"
	if destParent == `` {
		destParent = `/`
	}

	check(cli.CopyToContainer(context.Background(), c.ID, destParent, tar, types.CopyToContainerOptions{true, false}))
	log.Println(c.Names[0] + m.Destination)
}

//return nil if not found
func getContainerByName(cli *client.Client, targetName string) *types.Container {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	check(err)
	for _, curContainer := range containers {
		for _, curName := range curContainer.Names {
			if curName == `/`+targetName {
				return &curContainer
			}
		}
	}
	return nil
}

func restoreContainer(containerName string, wg *sync.WaitGroup) {
	defer wg.Done()

	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	c := getContainerByName(cli, containerName)
	if c == nil {
		log.Println(`Container "` + containerName + `" not found`)
		return
	}

	if c.State == `running` {
		if len(c.Mounts) > 0 {
			inspect, err := cli.ContainerInspect(context.Background(), c.ID)
			check(err)

			if inspect.HostConfig.AutoRemove == false {
				timeout := time.Minute
				check(cli.ContainerStop(context.Background(), c.ID, &timeout))

				var wgMount sync.WaitGroup
				wgMount.Add(len(c.Mounts))
				for _, curMount := range c.Mounts {
					go restoreMount(*c, curMount, &wgMount)
				}
				wgMount.Wait()

				check(cli.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{}))
			}
		}
	}
}
