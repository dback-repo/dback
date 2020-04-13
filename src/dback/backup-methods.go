package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

func packWithRestic(containerName, mountDestination, s3Endpoint, s3Bucket, accKey, secKey string) {

	cmd := exec.Command(`/bin/restic`, `init`)
	cmd.Dir = containerName
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
	//log.Printf("%s\n", stdoutStderr)

	files, err := ioutil.ReadDir(containerName)
	if err != nil {
		//log.Fatal(err)
		//panic(`sdf`)
	}

	log.Println(`===`)
	for _, file := range files {
		log.Println(`===`, file.Name())
	}

	log.Println(`----`, containerName+mountDestination)
	cmd = exec.Command(`/bin/restic`, `backup`, containerName+mountDestination)
	log.Println(`***`, `/bin/restic`, `backup`, containerName+mountDestination)
	cmd.Dir = containerName
	cmd.Env = append(os.Environ(),
		`RESTIC_PASSWORD=sdf`,
		`RESTIC_REPOSITORY=s3:http://`+s3Endpoint+`/`+s3Bucket+containerName+mountDestination,
		`AWS_ACCESS_KEY_ID=`+accKey,
		`AWS_SECRET_ACCESS_KEY=`+secKey)

	//log.Println(`---`, `RESTIC_REPOSITORY=/dback-snapshots`+containerName+mountDestination)
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		log.Println(string(stdoutStderr))
	}
	log.Printf("%s\n", stdoutStderr)

}

func unpackTarToMyself(c types.Container, myContainer types.Container, m types.MountPoint, s3Endpoint, s3Bucket, accKey, secKey string) {
	//tmp := fmt.Sprintf("%d", time.Now().UnixNano())

	//log.Println(tmp)
	//check(os.MkdirAll(c.Names[0], 664))

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
	// log.Println(`***`, m.Destination)
	// log.Println(`###`, destParent)
	os.MkdirAll(c.Names[0]+destParent, 0664)

	check(cli.CopyToContainer(context.Background(), myContainer.ID, c.Names[0]+destParent, tar, types.CopyToContainerOptions{true, false}))
	packWithRestic(c.Names[0], m.Destination, s3Endpoint, s3Bucket, accKey, secKey)
}

func backupMount(cli *client.Client, c types.Container, m types.MountPoint, wg *sync.WaitGroup, s3Endpoint, s3Bucket, accKey, secKey string) {
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

	unpackTarToMyself(c, *getContainerByNameOrId(cli, myselfContainerID), m, s3Endpoint, s3Bucket, accKey, secKey)
}

func backupContainer(c types.Container, wg *sync.WaitGroup, excludePattern, s3Endpoint, s3Bucket, accKey, secKey string) {
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
						go backupMount(cli, c, curMount, &wgMount, s3Endpoint, s3Bucket, accKey, secKey)
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
