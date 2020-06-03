package logic

import (
	"context"
	"dback/utils/cli"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
)

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, "\r\n", err.Error())
	}
}

func getContainersForBackup(dockerWrapper *dockerwrapper.DockerWrapper, matchers []string) []types.Container {
	allContainers := dockerWrapper.GetAllContainers()
	matchContainers := []types.Container{}

	for _, curContainer := range allContainers {
		if len(curContainer.Mounts) == 0 {
			log.Println(`Ignore container: `, curContainer.Names[0], ` cause: container has no mounts`)
		}

		_, cntBytes, _ := dockerWrapper.Docker.ContainerInspectWithRaw(context.Background(), curContainer.ID, true)

		match := true // container will be selected for backup, if inspect json contains all matchers substrings

		for _, curMatcher := range matchers {
			if !strings.Contains(string(cntBytes), curMatcher) {
				log.Println(`Ignore container: `, curContainer.Names[0], ` cause: matcher not found`, curMatcher)

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

func Backup(dockerWrapper *dockerwrapper.DockerWrapper, dbackOpts cli.DbackOpts,
	resticWrapper *resticwrapper.ResticWrapper) {
	mounts := getMountsForBackup(dockerWrapper, dbackOpts.Matchers, dbackOpts.ExcludePatterns)

	if isMountsEmpty(mounts) {
		log.Println(`No mounts for backup. Check "matcher" and "exclude" command line flags.
Run "dback backup --help" for more info`)
		return
	}

	if dbackOpts.IsEmulation {
		log.Println()
		log.Println(`Emulation started`)
		printMounts(mounts)
		log.Println(`The mounts above will be backup, if run dback without --emulate (-e) flag`)

		return
	}

	dbackOpts.ThreadsCount = correctThreadsCount(dbackOpts.ThreadsCount, len(mounts))

	log.Println()
	log.Println(`Backup started`)

	startBackupMoment := time.Now()

	saveMountsToResticParallel(dockerWrapper, mounts, dbackOpts.ThreadsCount, resticWrapper)
	log.Println(`Backup finished for the mounts above, in ` + time.Since(startBackupMoment).String())
}

func saveMountsToResticParallel(dockerWrapper *dockerwrapper.DockerWrapper, mounts []dockerwrapper.Mount,
	threadsCount int, resticWrapper *resticwrapper.ResticWrapper) {
	wg := sync.WaitGroup{}
	wg.Add(threadsCount)

	mountsCh := make(chan dockerwrapper.Mount)

	for i := 0; i < threadsCount; i++ {
		go saveMountsWorker(dockerWrapper, mountsCh, &wg, resticWrapper)
	}

	for _, curMount := range mounts {
		mountsCh <- curMount
	}

	close(mountsCh)
	wg.Wait()
}

func saveMountsWorker(dockerWrapper *dockerwrapper.DockerWrapper, ch chan dockerwrapper.Mount,
	wg *sync.WaitGroup, resticWrapper *resticwrapper.ResticWrapper) {
	for {
		mount, more := <-ch

		if !more {
			break
		}

		copyMountToLocal(dockerWrapper, mount)
		log.Println(`Save to restic:`, mount.ContainerName+mount.MountDest)
		resticWrapper.Save(`/tmp/dback-data/mount-data`+mount.ContainerName+mount.MountDest,
			mount.ContainerName+mount.MountDest)
	}
	wg.Done()
}

func copyMountToLocal(dockerWrapper *dockerwrapper.DockerWrapper, mount dockerwrapper.Mount) {
	dockerWrapper.CopyFolderToTar(mount.ContainerID, mount.MountDest,
		`/tmp/dback-data/tarballs`+mount.ContainerName+mount.MountDest)
	check(os.MkdirAll(`/tmp/dback-data/mount-data`+mount.ContainerName+mount.MountDest, 0664), `cannot make folder`)
	dockerWrapper.CopyTarToFloder(`/tmp/dback-data/tarballs`+mount.ContainerName+mount.MountDest+`/tar.tar`,
		dockerWrapper.GetMyselfContainerID(), `/tmp/dback-data/mount-data`+mount.ContainerName+mount.MountDest)
}
