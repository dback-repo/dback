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
	dockerw *dockerwrapper.DockerWrapper, containerName string) []s3wrapper.S3Mount {
	log.Println(`ContainerName`, containerName)
	s3Mounts := s3.GetMounts(resticw, dockerw, containerName)

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

	log.Println(`localMounts:`, localMounts)
	log.Println(`s3Mounts:`, res)

	return res
}

func isS3MountsEmpty(mounts []s3wrapper.S3Mount) bool {
	return len(mounts) == 0
}

func printS3MountsList(mounts []s3wrapper.S3Mount) {
	log.Println(`Mounts list:`, mounts)
}

const RestoreContainerTwoArgs = 2
const RestoreContainerThreeArgs = 3

func Restore(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper, dockerw *dockerwrapper.DockerWrapper,
	dbackParams cli.DbackOpts, dbackArgs []string, spacetracker *spacetracker.SpaceTracker) {
	if len(dbackArgs) == 0 {
		restore(s3, resticw, dockerw, dbackParams, spacetracker)
	} else {
		switch dbackArgs[0] {
		case `container`:
			if len(dbackArgs) == RestoreContainerTwoArgs {
				restoreContainer(s3, resticw, dockerw, dbackParams, dbackArgs[1], spacetracker)
			}
			if len(dbackArgs) == RestoreContainerThreeArgs {
				restoreContainerNewName(s3, resticw, dockerw, dbackParams, dbackArgs, spacetracker)
			}
		case `mount`:
			log.Fatalln(`Error: Mounts restoring is not implemented yet`)
		}
	}
}

func restore(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, dbackOpts cli.DbackOpts, spacetracker *spacetracker.SpaceTracker) {
	s3MountsForRestore := getS3MountsForRestore(s3, resticw, dockerw, ``)

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
		dockerWrapper.GetContainerIDByName(mount.ContainerName), destParent(mount.Dest))
}
