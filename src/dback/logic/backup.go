package logic

import (
	"dback/utils/dockerwrapper"
	"log"
	"os"
	"sync"
)

func check(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, "\r\n", err.Error())
	}
}

func Backup(dockerWrapper *dockerwrapper.DockerWrapper, isEmulation dockerwrapper.EmulateFlag,
	excludePatterns []dockerwrapper.ExcludePattern, threadsCount int) {
	log.Println(`Backup started`)

	containers := dockerWrapper.GetAllContainers()
	containers = dockerWrapper.SelectRunningContainers(containers)
	containers = dockerWrapper.SelectNotTemporaryContainers(containers)

	mounts := dockerWrapper.GetMountsOfContainers(containers)

	saveMountsToResticParallel(dockerWrapper, mounts, threadsCount)

	log.Println(mounts)
}

func saveMountsToResticParallel(dockerWrapper *dockerwrapper.DockerWrapper, mounts []dockerwrapper.Mount,
	threadsCount int) {
	wg := sync.WaitGroup{}
	wg.Add(threadsCount)

	mountsCh := make(chan dockerwrapper.Mount)

	for i := 0; i < threadsCount; i++ {
		go saveMountsWorker(dockerWrapper, mountsCh, &wg)
	}

	for _, curMount := range mounts {
		mountsCh <- curMount
	}

	close(mountsCh)
	wg.Wait()
	log.Println(`finished`)
}

func saveMountsWorker(dockerWrapper *dockerwrapper.DockerWrapper, ch chan dockerwrapper.Mount, wg *sync.WaitGroup) {
	for {
		mount, more := <-ch

		if !more {
			break
		}

		copyMountToLocal(dockerWrapper, mount)

	}
	wg.Done()
}

func pwd() string {
	res, err := os.Getwd()
	check(err, `cannot get current directory`)

	return res
}

func copyMountToLocal(dockerWrapper *dockerwrapper.DockerWrapper, mount dockerwrapper.Mount) {
	dockerWrapper.CopyFolderToTar(mount.ContainerID, mount.MountDest,
		`dback-data/tarballs`+mount.ContainerName+mount.MountDest)
	check(os.MkdirAll(`dback-data/mount-data`+mount.ContainerName+mount.MountDest, 0664), `cannot make folder`)
	dockerWrapper.CopyTarToFloder(`dback-data/tarballs`+mount.ContainerName+mount.MountDest+`/tar.tar`,
		dockerWrapper.GetMyselfContainerID(), pwd()+`dback-data/mount-data`+mount.ContainerName+mount.MountDest)
}
