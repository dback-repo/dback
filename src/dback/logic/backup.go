package logic

import (
	"bytes"
	"context"
	"dback/utils/cli"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/spacetracker"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/antchfx/htmlquery"
	"github.com/docker/docker/api/types"
	"github.com/yosssi/gohtml"
	"vimagination.zapto.org/json2xml"
)

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, "\r\n", err.Error())
	}
}

//check a node exist at least once
func isNodeExistByXpath(xmlNode *html.Node, xpath string) bool {
	return len(htmlquery.Find(xmlNode, xpath)) > 0
}

func getXmlInspectByContainer(dockerWrapper *dockerwrapper.DockerWrapper, container types.Container) *html.Node {
	_, cntBytes, _ := dockerWrapper.Docker.ContainerInspectWithRaw(context.Background(), container.ID, true)

	buf := strings.Builder{}
	x := xml.NewEncoder(&buf)
	check(json2xml.Convert(json.NewDecoder(bytes.NewReader(cntBytes)), x), `cannot convert json to xml`)
	check(x.Flush(), `cannot flush xml encoder`)

	gohtml.Condense = true
	res, err := htmlquery.Parse(bytes.NewReader([]byte(buf.String())))
	check(err, `Cannot parse xml: `+buf.String())

	return res
}

func getContainersForBackup(dockerWrapper *dockerwrapper.DockerWrapper, matchers []string) []types.Container {
	allContainers := dockerWrapper.GetAllContainers()
	matchContainers := []types.Container{}

	for _, curContainer := range allContainers {
		if len(curContainer.Mounts) == 0 {
			log.Println(`Ignore container: `, dockerWrapper.GetCorrectContainerName(curContainer.Names),
				` cause: container has no mounts`)
		}

		xmlInspectNode := getXmlInspectByContainer(dockerWrapper, curContainer)
		match := true // container will be selected for backup, if inspect json contains all matchers substrings

		for _, curMatcher := range matchers {
			if !isNodeExistByXpath(xmlInspectNode, curMatcher) {
				log.Println(`Ignore container: `, dockerWrapper.GetCorrectContainerName(curContainer.Names),
					` cause: matcher not found`, curMatcher)

				match = false

				break
			}
		}

		if match {
			matchContainers = append(matchContainers, curContainer)
		}
	}

	return matchContainers
}

//0    => mountCount
//1    => 1
//9999 => mountCount
func correctThreadsCount(threadsCount int, mountCount int) int {
	if threadsCount == 0 || threadsCount > mountCount {
		threadsCount = mountCount
	}

	return threadsCount
}

func isMountsEmpty(mounts []dockerwrapper.Mount) bool {
	return len(mounts) == 0
}

func getMountsForBackup(dockerWrapper *dockerwrapper.DockerWrapper, matchers []string,
	excludePatterns []dockerwrapper.ExcludePattern) []dockerwrapper.Mount {
	containers := getContainersForBackup(dockerWrapper, matchers)
	mounts := dockerWrapper.GetMountsOfContainers(containers)
	mounts = dockerWrapper.ExcludeMountsByPattern(mounts, excludePatterns)

	return mounts
}

func printMounts(mounts []dockerwrapper.Mount) {
	for _, curMount := range mounts {
		log.Println(curMount.ContainerName + curMount.MountDest)
	}
}

func newTimestamp(moment time.Time) string {
	return moment.Format(`02.01.2006.15-04-05`)
}

func backupEmulation(mounts []dockerwrapper.Mount) {
	log.Println()
	log.Println(`Emulation started`)
	printMounts(mounts)
	log.Println(`The mounts above will be backup, if run dback without --emulate (-e) flag`)
}

func printEmptyMountsMess() {
	log.Println(`No mounts for backup. Check "matcher" and "exclude" command line flags.
Run "dback backup --help" for more info`)
}

//9.5213121s -> 9s
func secondsFormat(t time.Duration) string {
	tstr := t.String()

	if t > time.Second {
		tstr = tstr[:strings.Index(tstr, `.`)]
	}

	return tstr + `s`
}

func Backup(dockerWrapper *dockerwrapper.DockerWrapper, dbackOpts cli.DbackOpts,
	resticWrapper *resticwrapper.ResticWrapper, spacetracker *spacetracker.SpaceTracker) {
	mounts := getMountsForBackup(dockerWrapper, dbackOpts.Matchers, dbackOpts.ExcludePatterns)

	if isMountsEmpty(mounts) {
		printEmptyMountsMess()
		return
	}

	if dbackOpts.IsEmulation {
		backupEmulation(mounts)
		return
	}

	dbackOpts.ThreadsCount = correctThreadsCount(dbackOpts.ThreadsCount, len(mounts))

	startBackupMoment := time.Now()
	timestamp := newTimestamp(startBackupMoment)

	log.Println()
	log.Println(`Backup started. Timestamp = ` + timestamp)
	saveMountsToResticParallel(dockerWrapper, mounts, dbackOpts.ThreadsCount, resticWrapper, timestamp)
	spacetracker.PrintReport()
	log.Println(`Backup finished for the mounts above, in ` + secondsFormat(time.Since(startBackupMoment)))
}

func startStoppedContainers(dockerWrapper *dockerwrapper.DockerWrapper, stoppedContainers []string) {
	dockerWrapper.StartContainersByIDs(&stoppedContainers, true)
}

func saveMountsToResticParallel(dockerWrapper *dockerwrapper.DockerWrapper, mounts []dockerwrapper.Mount,
	threadsCount int, resticWrapper *resticwrapper.ResticWrapper, timestamp string) {
	stoppedContainers := dockerWrapper.SelectRunningContainersByIDs(dockerWrapper.GetContainerIDsOfMounts(mounts))
	dockerWrapper.StopContainersByIDs(stoppedContainers, true)

	wg := sync.WaitGroup{}
	wg.Add(threadsCount)

	mountsCh := make(chan dockerwrapper.Mount)

	var saveErr error

	for i := 0; i < threadsCount; i++ {
		go saveMountsWorker(dockerWrapper, mountsCh, &wg, resticWrapper, timestamp, &saveErr)
	}

	for _, curMount := range mounts {
		mountsCh <- curMount
	}

	close(mountsCh)
	wg.Wait()
	startStoppedContainers(dockerWrapper, stoppedContainers) // errors with details are already printed

	if saveErr != nil {
		log.Fatalln(`Cannot save a mount. Stopped containers are switched back to run`)
	}
}

func saveMountsWorker(dockerWrapper *dockerwrapper.DockerWrapper, ch chan dockerwrapper.Mount,
	wg *sync.WaitGroup, resticWrapper *resticwrapper.ResticWrapper, timestamp string, saveErr *error) {
	defer func() {
		if r := recover(); r != nil {
			*saveErr = errors.New(fmt.Sprint(r))
		}

		wg.Done()
	}()

	for {
		mount, more := <-ch

		if !more {
			break
		}

		copyMountToLocal(dockerWrapper, mount)

		log.Println(`Save to restic:`, mount.ContainerName+mount.MountDest)
		resticWrapper.Save(`/tmp/dback-data/mount-data`+mount.ContainerName+mount.MountDest,
			mount.ContainerName+mount.MountDest, timestamp)
	}
}

func copyMountToLocal(dockerWrapper *dockerwrapper.DockerWrapper, mount dockerwrapper.Mount) {
	dockerWrapper.CopyFolderToTar(mount.ContainerID, mount.MountDest,
		`/tmp/dback-data/tarballs`+mount.ContainerName+mount.MountDest)
	check(os.MkdirAll(destParent(`/tmp/dback-data/mount-data`+mount.ContainerName+
		mount.MountDest), 0664), `cannot make folder`)
	dockerWrapper.CopyTarToFloder(`/tmp/dback-data/tarballs`+mount.ContainerName+mount.MountDest+`/tar.tar`,
		dockerWrapper.GetMyselfContainerID(), destParent(`/tmp/dback-data/mount-data`+mount.ContainerName+mount.MountDest))

	go check(os.RemoveAll(`/tmp/dback-data/tarballs`+mount.ContainerName+mount.MountDest+`/tar.tar`), `cannot remove tar`)
}

func destParent(dest string) string {
	lastSlashIdx := strings.LastIndex(dest, `/`)
	destParent := dest[:lastSlashIdx] //      "/var/www/lynx" -> "/var/www"        "/opt" -> "/"

	if destParent == `` {
		destParent = `/`
	}

	return destParent
}
