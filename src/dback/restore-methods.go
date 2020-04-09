package main

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
)

//list of saved containers == list of folders in /dback-snapshots
func getBackupsContainerList(s3Endpoint, s3Bucket, accKey, secKey string) []string {
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(s3Endpoint, accKey, secKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	doneCh := make(chan struct{})
	defer close(doneCh)

	containers := []string{}
	for object := range minioClient.ListObjects(s3Bucket, ``, false, doneCh) {
		check(object.Err)
		containers = append(containers, object.Key)
	}

	return containers
}

func restoreContainers(containers []string) {
	log.Println(containers)
	return

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

	tar, err := os.Open(`dback-snapshots/` + c.Names[0] + m.Destination + `/tar.tar`)
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
func getContainerByNameOrId(cli *client.Client, targetName string) *types.Container {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	check(err)
	for _, curContainer := range containers {
		for _, curName := range curContainer.Names {
			if curName == `/`+targetName {
				return &curContainer
			}
		}
		//log.Println(targetName, ` `, curContainer.ID[:len(targetName)])
		if curContainer.ID[:len(targetName)] == targetName {
			return &curContainer
		}
	}
	return nil
}

func restoreContainer(containerName string, wg *sync.WaitGroup) {
	defer wg.Done()

	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	c := getContainerByNameOrId(cli, containerName)
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
