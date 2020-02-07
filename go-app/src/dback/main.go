package main

import (
	"context"
	"io"

	//	"log"
	"os"
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

	check(os.MkdirAll(`/backup/`+c.Names[0]+`/`+m.Destination, 0664))

	reader, _, err := cli.CopyFromContainer(context.Background(), c.ID, m.Destination)
	check(err)

	outFile, err := os.Create(`/backup/` + c.Names[0] + `/` + m.Destination + `/mount-data.tar`)
	check(err)
	defer outFile.Close()
	_, err = io.Copy(outFile, reader)
	check(err)
}

func backupContainer(c types.Container, wg *sync.WaitGroup) {
	defer wg.Done()

	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	timeout := time.Minute
	check(cli.ContainerStop(context.Background(), c.ID, &timeout))

	var wgMount sync.WaitGroup
	wgMount.Add(len(c.Mounts))

	for _, curMount := range c.Mounts {
		go backupMount(c, curMount, &wgMount)
	}

	wgMount.Wait()
}

func main() {
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
}
