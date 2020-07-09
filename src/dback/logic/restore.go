package logic

import (
	"dback/utils/cli"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"dback/utils/spacetracker"
	"log"
	"os"
	"strings"
	"sync"
)

func getS3MountsOfContainer(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, containerName string) []s3wrapper.S3Mount {
	containerName = strings.TrimPrefix(containerName, `/`)

	res := s3.GetMounts(resticw, dockerw, containerName)

	for mountIdx, curS3Mount := range res {
		res[mountIdx].SelectedSnapshotID = curS3Mount.Snapshots[len(curS3Mount.Snapshots)-1].ID
	}

	return res
}

func getS3MountsForRestore(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, prefix string) []s3wrapper.S3Mount {
	prefix = strings.TrimPrefix(prefix, `/`)
	s3Mounts := s3.GetMounts(resticw, dockerw, prefix)

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

	return res
}

func isS3MountsEmpty(mounts []s3wrapper.S3Mount) bool {
	return len(mounts) == 0
}

const RestoreContainerTwoArgs = 2
const RestoreContainerThreeArgs = 3

func Restore(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper, dockerw *dockerwrapper.DockerWrapper,
	dbackParams cli.DbackOpts, dbackArgs []string, spacetracker *spacetracker.SpaceTracker) {
	if len(dbackArgs) == 0 {
		restoreAll(s3, resticw, dockerw, dbackParams, spacetracker)
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
			if len(dbackArgs) == RestoreContainerTwoArgs {
				restoreMount(s3, resticw, dockerw, dbackArgs[1], dbackArgs[1], dbackParams.Snapshot, spacetracker)
			} else {
				restoreMount(s3, resticw, dockerw, dbackArgs[1], dbackArgs[2], dbackParams.Snapshot, spacetracker)
			}
		}
	}
}

func restoreAll(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, dbackOpts cli.DbackOpts, spacetracker *spacetracker.SpaceTracker) {
	s3MountsForRestore := getS3MountsForRestore(s3, resticw, dockerw, ``)

	if isS3MountsEmpty(s3MountsForRestore) {
		printEmptyMountsMess()
		return
	}

	dbackOpts.ThreadsCount = correctThreadsCount(dbackOpts.ThreadsCount, len(s3MountsForRestore))

	loadMountsFromResticParallel(dockerw, s3MountsForRestore, dbackOpts.ThreadsCount, dbackOpts.Snapshot, resticw)
	spacetracker.PrintReport()
	log.Println(`Restoring finished for the mounts above`)
}

func loadMountsFromResticParallel(dockerWrapper *dockerwrapper.DockerWrapper, s3Mounts []s3wrapper.S3Mount,
	threadsCount int, snapshot string, resticWrapper *resticwrapper.ResticWrapper) {
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
		curMount.SelectSnapshotByTag(snapshot)

		if curMount.SelectedSnapshotID == `` {
			log.Fatalln(`Snapshot `, snapshot, ` not found for mount`, curMount.ContainerName+curMount.Dest)
		}

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
		copyLocalToMount(dockerWrapper, mount, ``)

		log.Println(`Load from restic:`, mount.ContainerName+mount.Dest)
	}
	wg.Done()
}

func copyLocalToMount(dockerWrapper *dockerwrapper.DockerWrapper, mount s3wrapper.S3Mount, containerName string) {
	if containerName == `` {
		containerName = mount.ContainerName
	}

	dockerWrapper.CopyFolderToTar(dockerWrapper.GetMyselfContainerID(), `/tmp/dback-data/mount-data`+
		mount.ContainerName+mount.Dest, `/tmp/dback-data/tarballs`+mount.ContainerName+mount.Dest)
	dockerWrapper.CopyTarToFloder(`/tmp/dback-data/tarballs`+mount.ContainerName+mount.Dest+`/tar.tar`,
		dockerWrapper.GetContainerIDByName(containerName), destParent(mount.Dest))
}
