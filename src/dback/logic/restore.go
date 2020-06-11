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

func getS3MountsForRestore(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, restoreParams []string) []s3wrapper.S3Mount {
	s3Mounts := s3.GetMounts(resticw, dockerw)

	localMounts := dockerw.GetMountsOfContainers(dockerw.GetAllContainers())

	res := []s3wrapper.S3Mount{}

	for _, curLocalMount := range localMounts {
		for _, curS3Mount := range s3Mounts {
			if curS3Mount.ContainerName+curS3Mount.Dest == curLocalMount.ContainerName+curLocalMount.MountDest {
				res = append(res, curS3Mount)
			}
		}
	}

	for mountIdx, curS3Mount := range res {
		res[mountIdx].SelectedSnapshotID = curS3Mount.Snapshots[len(curS3Mount.Snapshots)-1].ID
	}

	//log.Println(`localMounts:`, localMounts)
	log.Println(`restoreParams:`, restoreParams)

	return res
}

func isS3MountsEmpty(mounts []s3wrapper.S3Mount) bool {
	return len(mounts) == 0
}

func printS3MountsList(mounts []s3wrapper.S3Mount) {
	log.Println(`Mounts list:`, mounts)
}

func Restore(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, dbackOpts cli.DbackOpts, spacetracker *spacetracker.SpaceTracker) {
	s3MountsForRestore := getS3MountsForRestore(s3, resticw, dockerw, dbackOpts.Matchers)

	if isS3MountsEmpty(s3MountsForRestore) {
		printEmptyMountsMess()
		return
	}

	printS3MountsList(s3MountsForRestore)

	// if !isApproved(dbackOpts.AutoProceedFlag) {
	// 	return
	// }

	dbackOpts.ThreadsCount = correctThreadsCount(dbackOpts.ThreadsCount, len(s3MountsForRestore))

	loadMountsFromResticParallel(dockerw, s3MountsForRestore, dbackOpts.ThreadsCount, resticw)
	spacetracker.PrintReport()
	log.Println(`Restore finished for the mounts above`)
}

func loadMountsFromResticParallel(dockerWrapper *dockerwrapper.DockerWrapper, s3Mounts []s3wrapper.S3Mount,
	threadsCount int, resticWrapper *resticwrapper.ResticWrapper) {
	stoppedContainers := []string{}
	defer dockerWrapper.StartContainersByIDs(&stoppedContainers, false) // for start containers even after panic

	stoppedContainers = dockerWrapper.SelectRunningContainersByIDs(s3wrapper.GetContainerIDsOfMounts(
		s3Mounts, dockerWrapper))
	dockerWrapper.StopContainersByIDs(stoppedContainers, true)

	wg := sync.WaitGroup{}
	wg.Add(threadsCount)

	mountsCh := make(chan s3wrapper.S3Mount)

	for i := 0; i < threadsCount; i++ {
		go loadMountsWorker(dockerWrapper, mountsCh, &wg, resticWrapper)
	}

	for _, curMount := range s3Mounts {
		mountsCh <- curMount
	}

	close(mountsCh)
	wg.Wait()
}

func loadMountsWorker(dockerWrapper *dockerwrapper.DockerWrapper, ch chan s3wrapper.S3Mount,
	wg *sync.WaitGroup, resticWrapper *resticwrapper.ResticWrapper) {
	for {
		mount, more := <-ch

		if !more {
			break
		}

		check(os.MkdirAll(`/tmp/dback-data/mount-data`+mount.ContainerName+mount.Dest, 0664), `cannot make folder`)
		resticWrapper.Load(`/`, mount.ContainerName+mount.Dest, mount.SelectedSnapshotID)
		copyLocalToMount(dockerWrapper, mount)

		log.Println(`Load from restic:`, mount.ContainerName+mount.Dest)
	}
	wg.Done()
}

func copyLocalToMount(dockerWrapper *dockerwrapper.DockerWrapper, mount s3wrapper.S3Mount) {
	dockerWrapper.CopyFolderToTar(dockerWrapper.GetMyselfContainerID(), `/tmp/dback-data/mount-data`+
		mount.ContainerName+mount.Dest, `/tmp/dback-data/tarballs`+mount.ContainerName+mount.Dest)
	dockerWrapper.CopyTarToFloder(`/tmp/dback-data/tarballs`+mount.ContainerName+mount.Dest+`/tar.tar`,
		dockerWrapper.GetContainerIDByName(mount.ContainerName), `/`)
}
