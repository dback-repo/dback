package logic

import (
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"dback/utils/spacetracker"
	"log"
	"os"
	"strings"
)

func contanerNameByS3MountID(s3Mount string) string {
	res := ``
	arr := strings.Split(s3Mount, `/`)

	if len(arr) >= 1 {
		res = arr[1]
	}

	return res
}

func restoreMount(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, mount1, mount2, snapshot string,
	spacetracker *spacetracker.SpaceTracker) {
	s3MountsForRestore := getS3MountsForRestore(s3, resticw, dockerw, mount1)

	if isS3MountsEmpty(s3MountsForRestore) {
		printEmptyMountsMess()
		return
	}

	var s3Mount *s3wrapper.S3Mount

	for curMountIdx, curMount := range s3MountsForRestore {
		if curMount.ContainerName+curMount.Dest == mount1 {
			s3Mount = &(s3MountsForRestore[curMountIdx])
			s3Mount.SelectSnapshotByTag(snapshot)

			break
		}
	}

	if s3Mount.SelectedSnapshotID == `` {
		log.Println(`Snapshot `+snapshot+` not found for mount`, s3Mount.ContainerName+s3Mount.Dest)
		return
	}

	if s3Mount == nil {
		printEmptyMountsMess()
		return
	}

	log.Println(`Load from restic: ` + mount1 + ` to ` + mount2)

	check(os.MkdirAll(`/tmp/dback-data/mount-data`+mount1, 0664), `cannot make folder`)
	resticw.Load(`/`, mount1, s3Mount.SelectedSnapshotID)
	copyLocalToMount(dockerw, *s3Mount, containerNameLeadingSlash(contanerNameByS3MountID(mount2)))

	spacetracker.PrintReport()
	log.Println(`Restore finished for the mounts above`)
}
