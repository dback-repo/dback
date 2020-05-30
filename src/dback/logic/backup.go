package logic

import (
	"dback/utils/cli"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"log"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
)

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, "\r\n", err.Error())
	}
}

func getContainersForBackup(dockerWrapper *dockerwrapper.DockerWrapper) []types.Container {
	containers := dockerWrapper.GetAllContainers()
	containers = dockerWrapper.SelectRunningContainers(containers)
	containers = dockerWrapper.SelectNotTemporaryContainers(containers)

	return containers
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
	if len(mounts) == 0 {
		log.Println(`No mounts for backup. Check "matcher" and "exclude" command line flags.
Run "dback backup --help" for more info`)

		return true
	}

	return false
}

func Backup(dockerWrapper *dockerwrapper.DockerWrapper, dbackOpts cli.DbackOpts,
	resticWrapper *resticwrapper.ResticWrapper) {
	containers := getContainersForBackup(dockerWrapper)
	mounts := dockerWrapper.GetMountsOfContainers(containers)

	if isMountsEmpty(mounts) {
		return
	}

	dbackOpts.ThreadsCount = correctThreadsCount(dbackOpts.ThreadsCount, len(mounts))

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
	log.Println(mount.ContainerName + mount.MountDest)
}
