package main

import (
	"context"
	"io"
	"strings"

	//"io/ioutil"
	"log"
	"os"
	"os/exec"

	//	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
)

func pullFilesFromRestic(containerName, mountDestination, s3Endpoint, s3Bucket, accKey, secKey string) {

	cmd := exec.Command(`/bin/restic`, `restore`, `latest`, `--target`, `.`)
	cmd.Dir = `/`
	//cmd.Env = append(os.Environ(), `RESTIC_REPOSITORY=/dback-snapshots`+containerName+mountDestination, `RESTIC_PASSWORD=sdf`)
	cmd.Env = append(os.Environ(),
		`RESTIC_PASSWORD=sdf`,
		`RESTIC_REPOSITORY=s3:http://`+s3Endpoint+`/`+s3Bucket+containerName+mountDestination,
		`AWS_ACCESS_KEY_ID=`+accKey,
		`AWS_SECRET_ACCESS_KEY=`+secKey)
	//s3:https://s3.amazonaws.com/BUCKET_NAME
	//log.Println(`---`, `RESTIC_REPOSITORY=/dback-snapshots`+containerName+mountDestination)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		log.Println(string(stdoutStderr))
		//panic(`sdf`)
	}
	log.Printf("%s\n", stdoutStderr)

	// files, err := ioutil.ReadDir(`/`)
	// check(err)
	// for _, f := range files {
	// 	fmt.Println(`->`, f.Name())
	// }

	// files, err := ioutil.ReadDir(tmp)
	// if err != nil {
	// 	//log.Fatal(err)
	// 	//panic(`sdf`)
	// }

	// log.Println(`===`)
	// for _, file := range files {
	// 	log.Println(`===`, file.Name())
	// }

	// log.Println(`----`, tmp+mountDestination)
	// cmd = exec.Command(`/bin/restic`, `backup`, tmp+mountDestination)
	// log.Println(`***`, `/bin/restic`, `backup`, tmp+mountDestination)
	// cmd.Dir = tmp
	// cmd.Env = append(os.Environ(),
	// 	`RESTIC_PASSWORD=sdf`,
	// 	`RESTIC_REPOSITORY=s3:http://`+s3Endpoint+`/`+s3Bucket+containerName+mountDestination,
	// 	`AWS_ACCESS_KEY_ID=`+accKey,
	// 	`AWS_SECRET_ACCESS_KEY=`+secKey)

	// //log.Println(`---`, `RESTIC_REPOSITORY=/dback-snapshots`+containerName+mountDestination)
	// stdoutStderr, err = cmd.CombinedOutput()
	// if err != nil {
	// 	log.Println(err)
	// 	log.Println(string(stdoutStderr))
	// }
	// log.Printf("%s\n", stdoutStderr)
}

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
		containers = append(containers, object.Key[:len(object.Key)-1])
	}

	return containers
}

func restoreContainers(containers []string, s3Endpoint, s3Bucket, accKey, secKey string) {
	// log.Println(containers)
	// return

	var wg sync.WaitGroup
	wg.Add(len(containers))

	for _, curContainer := range containers {
		go restoreContainer(curContainer, &wg, s3Endpoint, s3Bucket, accKey, secKey)
	}

	wg.Wait()
}

func packTarFromMyself(c types.Container, myContainer types.Container, m types.MountPoint) {
	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	check(os.MkdirAll(`dback-snapshots/`+c.Names[0]+m.Destination, 0664))

	reader, _, err := cli.CopyFromContainer(context.Background(), myContainer.ID, c.Names[0]+m.Destination)
	check(err)

	outFile, err := os.Create(`dback-snapshots/` + c.Names[0] + m.Destination + `/tar.tar`)
	check(err)
	defer outFile.Close()
	_, err = io.Copy(outFile, reader)
	check(err)

	//log.Println(`make backup: ` + c.Names[0] + m.Destination)

	// myselfContainerID, err := os.Hostname()
	// check(err)
}

func restoreMount(c types.Container, m types.MountPoint, wg *sync.WaitGroup, s3Endpoint, s3Bucket, accKey, secKey string) {
	defer wg.Done()

	cli, err := client.NewEnvClient()
	check(err)
	defer cli.Close()

	pullFilesFromRestic(c.Names[0], m.Destination, s3Endpoint, s3Bucket, accKey, secKey)

	myselfContainerID, err := os.Hostname()
	check(err)

	packTarFromMyself(c, *getContainerByNameOrId(cli, myselfContainerID), m)

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

func restoreContainer(containerName string, wg *sync.WaitGroup, s3Endpoint, s3Bucket, accKey, secKey string) {
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
					go restoreMount(*c, curMount, &wgMount, s3Endpoint, s3Bucket, accKey, secKey)
				}
				wgMount.Wait()

				check(cli.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{}))
			}
		}
	}
}
