package main

import (
	"context"
	"io/ioutil"
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

	check(cli.CopyToContainer(context.Background(), c.ID, ``, rdr, types.CopyToContainerOptions{true, false}))

	// 	check(os.MkdirAll(`/backup/`+c.Names[0]+m.Destination, 0664))

	// 	reader, _, err := cli.CopyFromContainer(context.Background(), c.ID, m.Destination)
	// 	check(err)

	// 	lastSlashIdx := strings.LastIndex(m.Destination, `/`)
	// 	// if lastSlashIdx > 0 {
	// 	// 	lastSlashIdx--
	// 	// }
	// 	log.Println(`lastSlashIdx`, lastSlashIdx)
	// 	destParent := m.Destination[:lastSlashIdx] // /var/www/lynx -> /var/www
	// 	if destParent == `` {
	// 		destParent = `/`
	// 	}
	// 	log.Println(`dest`, m.Destination)
	// 	log.Println(`destParent`, destParent)

	// 	check(Untar(reader, `/backup/`+c.Names[0]+destParent))
	// 	log.Println(c.Names[0] + m.Destination)
}

//return nil if not found
func getContainerByName(cli *client.Client, targetName string) *types.Container {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	check(err)
	for _, curContainer := range containers {
		for _, curName := range curContainer.Names {
			if curName == targetName {
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
