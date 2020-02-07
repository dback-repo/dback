package main

import (
	"context"
	"io"
	"log"
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

func backupContainer(c types.Container, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(c.Mounts)
	log.Println(c.State)
	log.Println(c.Status)
	cli, err := client.NewEnvClient()
	check(err)

	timeout := time.Minute
	check(cli.ContainerStop(context.Background(), c.ID, &timeout))
	log.Println(c.State)
	log.Println(c.Status)

	reader, _, err := cli.CopyFromContainer(context.Background(), c.ID, `/mount`)
	check(err)

	outFile, err := os.Create(`/backup/tar.tar`)
	check(err)
	defer outFile.Close()
	_, err = io.Copy(outFile, reader)
	//reader.
	//types.ContainerPathStat
}

func main() {
	cli, err := client.NewEnvClient()
	check(err)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	check(err)

	var wg sync.WaitGroup
	wg.Add(len(containers))

	for _, curContainer := range containers {
		go backupContainer(curContainer, &wg)
	}

	wg.Wait()
}
