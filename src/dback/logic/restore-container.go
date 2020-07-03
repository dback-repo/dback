package logic

import (
	"dback/utils/cli"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"dback/utils/spacetracker"
	"log"
	"os"
	"sync"
)

// // /lynx -> /lynx
// //  lynx -> /lynx
func containerNameLeadingSlash(name string) string {
	if len(name) > 0 {
		if string(name[0]) != `/` {
			name = `/` + name
		}
	}

	return name
}

func restoreContainer(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, dbackOpts cli.DbackOpts, containerName string,
	spacetracker *spacetracker.SpaceTracker) {
	s3MountsForRestore := getS3MountsForRestore(s3, resticw, dockerw, containerName)

	if isS3MountsEmpty(s3MountsForRestore) {
		printEmptyMountsMess()
		return
	}

	dbackOpts.ThreadsCount = correctThreadsCount(dbackOpts.ThreadsCount, len(s3MountsForRestore))

	loadMountsFromResticParallel(dockerw, s3MountsForRestore, dbackOpts.ThreadsCount, resticw)
	spacetracker.PrintReport()
	log.Println(`Restore finished for the mounts above`)
}

func restoreContainerNewName(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, dbackOpts cli.DbackOpts, dbackArgs []string,
	spacetracker *spacetracker.SpaceTracker) {
	s3MountsForRestore := getS3MountsForRestore(s3, resticw, dockerw, dbackArgs[1])

	if isS3MountsEmpty(s3MountsForRestore) {
		printEmptyMountsMess()
		return
	}

	dbackOpts.ThreadsCount = correctThreadsCount(dbackOpts.ThreadsCount, len(s3MountsForRestore))

	loadMountsFromResticParallelToContainer(dockerw, s3MountsForRestore, dbackOpts.ThreadsCount, resticw, dbackArgs[2])
	spacetracker.PrintReport()
	log.Println(`Restore finished for the mounts above`)
}

func loadMountsFromResticParallelToContainer(dockerWrapper *dockerwrapper.DockerWrapper,
	s3Mounts []s3wrapper.S3Mount, threadsCount int, resticWrapper *resticwrapper.ResticWrapper, containerName string) {
	stoppedContainers := []string{}
	defer dockerWrapper.StartContainersByIDs(&stoppedContainers, false) // for start containers even after panic

	stoppedContainers = dockerWrapper.SelectRunningContainersByIDs(s3wrapper.GetContainerIDsOfMounts(
		s3Mounts, dockerWrapper))
	dockerWrapper.StopContainersByIDs(stoppedContainers, true)

	wg := sync.WaitGroup{}
	wg.Add(threadsCount)

	mountsCh := make(chan s3wrapper.S3Mount)

	for i := 0; i < threadsCount; i++ {
		go loadMountsWorkerToContainer(dockerWrapper, mountsCh, &wg, resticWrapper, containerName)
	}

	for _, curMount := range s3Mounts {
		mountsCh <- curMount
	}

	close(mountsCh)
	wg.Wait()
}

func loadMountsWorkerToContainer(dockerWrapper *dockerwrapper.DockerWrapper, ch chan s3wrapper.S3Mount,
	wg *sync.WaitGroup, resticWrapper *resticwrapper.ResticWrapper, containerName string) {
	for {
		mount, more := <-ch

		if !more {
			break
		}

		log.Println(`Load from restic: ` + mount.ContainerName + mount.Dest + ` to /` +
			containerName + mount.Dest)

		check(os.MkdirAll(`/tmp/dback-data/mount-data`+containerNameLeadingSlash(containerName)+
			mount.Dest, 0664), `cannot make folder`)
		resticWrapper.Load(`/`, mount.ContainerName+mount.Dest, mount.SelectedSnapshotID)
		copyLocalToMount(dockerWrapper, mount, containerNameLeadingSlash(containerName))
	}
	wg.Done()
}
