package dockerwrapper

import (
	"context"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

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

func checkJustLog(err error, msg string) {
	if err != nil {
		log.Println(msg + "\r\n" + err.Error())
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

func (t *DockerWrapper) GetMountsOfContainers(containers []types.Container) []Mount {
	res := []Mount{}

	for _, curContainer := range containers {
		for _, curMount := range curContainer.Mounts {
			res = append(res, Mount{curContainer.ID, t.GetCorrectContainerName(curContainer.Names), curMount.Destination})
		}
	}

	return res
}

func (t *DockerWrapper) ExcludeMountsByPattern(mounts []Mount, excludePatterns []ExcludePattern) []Mount {
	mountsForBackup := []Mount{}

	for _, curMount := range mounts {
		backupMount := true

		for _, curExcludePattern := range excludePatterns {
			r, err := regexp.Compile(string(curExcludePattern))
			check(err, `Exclude pattern is not correct regexp`+string(curExcludePattern))

			if r.MatchString(curMount.ContainerName + curMount.MountDest) {
				backupMount = false

				log.Println(`Exclude mount: ` + curMount.ContainerName + curMount.MountDest +
					`      cause: --exclude-mount ` + string(curExcludePattern))
			}
		}

		if backupMount {
			mountsForBackup = append(mountsForBackup, curMount)
		}
	}

	return mountsForBackup
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

func (t *DockerWrapper) GetContainerIDByName(containerName string) string {
	containers := t.GetAllContainers()
	for _, curContainer := range containers {
		if t.GetCorrectContainerName(curContainer.Names) == containerName {
			return curContainer.ID
		}
	}

	return ``
}

func (t *DockerWrapper) GetCorrectContainerName(names []string) string {
	for _, curName := range names {
		if strings.Count(curName, `/`) == 1 {
			return curName
		}
	}

	return ``
}

func (t *DockerWrapper) GetContainerIDsOfMounts(mounts []Mount) []string {
	res := []string{}

	containersMap := make(map[string]string)

	for _, curMount := range mounts {
		containersMap[curMount.ContainerID] = curMount.ContainerID
	}

	for _, curContainerID := range containersMap {
		res = append(res, curContainerID)
	}

	return res
}

func (t *DockerWrapper) StartContainersByIDs(ids *[]string, panicOnError bool) {
	for _, curContainerID := range *ids {
		err := t.Docker.ContainerStart(context.Background(), curContainerID, types.ContainerStartOptions{})

		log.Println(`StartContainer:`, curContainerID)

		if panicOnError {
			check(err, `Cannot stop container: `+curContainerID)
		} else {
			checkJustLog(err, `Cannot stop container: `+curContainerID)
		}
	}
}

func (t *DockerWrapper) StopContainersByIDs(ids []string, panicOnError bool) {
	timeout := time.Minute

	for _, curContainerID := range ids {
		err := t.Docker.ContainerStop(context.Background(), curContainerID, &timeout)

		log.Println(`StopContainer:`, curContainerID)

		if panicOnError {
			check(err, `Cannot stop container: `+curContainerID)
		} else {
			checkJustLog(err, `Cannot stop container: `+curContainerID)
		}
	}
}

func (t *DockerWrapper) SelectRunningContainersByIDs(ids []string) []string {
	res := []string{}

	for _, curContainerID := range ids {
		inspect, err := t.Docker.ContainerInspect(context.Background(), curContainerID)
		check(err, `cannot inspect container`+curContainerID)

		if inspect.State.Running {
			res = append(res, curContainerID)
		}
	}

	return res
}
