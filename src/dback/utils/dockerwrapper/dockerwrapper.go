package dockerwrapper

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerWrapper struct {
	Docker *client.Client
}

type Mount struct {
	ContainerID   string
	ContainerName string
	MountDest     string
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg + "\r\n" + err.Error())
	}
}

func (t *DockerWrapper) Close() {
	check(t.Docker.Close(), `Cannot close docker connection`)
}

func (t *DockerWrapper) GetAllContainers() []types.Container {
	containers, err := t.Docker.ContainerList(context.Background(), types.ContainerListOptions{})
	check(err, `cannot get list of containers`)

	return containers
}

func (t *DockerWrapper) SelectRunningContainers(containers []types.Container) []types.Container {
	res := []types.Container{}

	for _, curContainer := range containers {
		if curContainer.State == `running` {
			res = append(res, curContainer)
		}
	}

	return res
}

func (t *DockerWrapper) SelectNotTemporaryContainers(containers []types.Container) []types.Container {
	res := []types.Container{}

	for _, curContainer := range containers {
		inspect, err := t.Docker.ContainerInspect(context.Background(), curContainer.ID)
		check(err, `Cannot inspect container `+curContainer.ID)

		if !inspect.HostConfig.AutoRemove {
			res = append(res, curContainer)
		}
	}

	return res
}

func (t *DockerWrapper) GetMountsOfContainers(containers []types.Container) []Mount {
	res := []Mount{}

	for _, curContainer := range containers {
		for _, curMount := range curContainer.Mounts {
			res = append(res, Mount{curContainer.ID, curContainer.Names[0], curMount.Destination})
		}
	}

	return res
}

func (t *DockerWrapper) CopyFolderToTar(containerID, folderDestination, tarDestination string) {
	check(os.MkdirAll(tarDestination, 0664), `cannot make directory `+tarDestination)

	reader, _, err := t.Docker.CopyFromContainer(context.Background(), containerID, folderDestination)
	check(err, `cannot copy from container `+containerID)

	outFile, err := os.Create(tarDestination + `/tar.tar`)
	check(err, `cannot create output file `+tarDestination+`/tar.tar`)

	defer outFile.Close()

	_, err = io.Copy(outFile, reader)
	check(err, `cannot copy io flow`)
}

//it is hostname
func (t *DockerWrapper) GetMyselfContainerID() string {
	res, err := os.Hostname()
	check(err, `cannot lookup a hostname`)

	return res
}

func (t *DockerWrapper) CopyTarToFloder(tarDestination, containerID, folderDestination string) {
	tar, err := os.Open(tarDestination)
	check(err, `cannot open file `+tarDestination)

	defer tar.Close()

	check(t.Docker.CopyToContainer(context.Background(), containerID,
		folderDestination, tar,
		types.CopyToContainerOptions{AllowOverwriteDirWithFile: true, CopyUIDGID: false}), `cannot copy to container`)
}
