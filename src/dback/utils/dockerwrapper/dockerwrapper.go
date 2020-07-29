package dockerwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/yosssi/gohtml"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"vimagination.zapto.org/json2xml"
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
	containers, err := t.Docker.ContainerList(context.Background(), types.ContainerListOptions{All: true})
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
		types.CopyToContainerOptions{AllowOverwriteDirWithFile: true, CopyUIDGID: true}), `cannot copy to container`)
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

//try to start all containers in the list
//return containers failed to start, and errors
func (t *DockerWrapper) startContainersByIDs(ids *[]string) (*[]string, []error) {
	res := []string{}
	errors := []error{}

	for _, curContainerID := range *ids {
		err := t.Docker.ContainerStart(context.Background(), curContainerID, types.ContainerStartOptions{})

		if err != nil {
			errors = append(errors, err)
			res = append(res, curContainerID)
		}
	}

	return &res, errors
}

const StartRetries = 10

func printContainersStartingErrors(errors []error, description string, panicOnError bool) {
	for _, curError := range errors {
		log.Println(curError.Error())
	}

	if panicOnError {
		log.Fatalln(description)
	}
}

func (t *DockerWrapper) StartContainersByIDs(ids *[]string, panicOnError bool) {
	var errors []error

	for i := 0; i < StartRetries; i++ {
		if len(*ids) == 0 {
			return
		}

		containersCount := len(*ids)
		ids, errors = t.startContainersByIDs(ids)

		if containersCount == len(*ids) { // if last call don't start any containers
			printContainersStartingErrors(errors, `Multiple errors while containers starting`, panicOnError)
		}
	}
}

func (t *DockerWrapper) StopContainersByIDs(ids []string, panicOnError bool) {
	timeout := time.Minute

	for _, curContainerID := range ids {
		err := t.Docker.ContainerStop(context.Background(), curContainerID, &timeout)

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

func (t *DockerWrapper) GetContainerNameByID(containerID string) string {
	for _, curContainer := range t.GetAllContainers() {
		if curContainer.ID == containerID {
			return t.GetCorrectContainerName(curContainer.Names)
		}
	}

	return ``
}

//it is hostname
func (t *DockerWrapper) GetDockerInspectXMLByContainerName(containerName string) string {
	_, cntBytes, err := t.Docker.ContainerInspectWithRaw(context.Background(), t.GetContainerIDByName(containerName), true)
	check(err, `cannot inspect container`)

	buf := strings.Builder{}
	x := xml.NewEncoder(&buf)
	check(json2xml.Convert(json.NewDecoder(bytes.NewReader(cntBytes)), x), `cannot convert json to xml`)
	check(x.Flush(), `cannot flush xml encoder`)

	gohtml.Condense = true

	return gohtml.Format(buf.String())
}
