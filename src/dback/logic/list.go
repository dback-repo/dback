package logic

import (
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"dback/utils/spacetracker"
	"log"
	"sort"
	"strings"
)

func printSortedS3Mounts(mounts []s3wrapper.S3Mount) {
	mountsStr := []string{}

	for _, curMount := range mounts {
		mountsStr = append(mountsStr, curMount.ContainerName+curMount.Dest)
	}

	sort.Strings(mountsStr)

	for _, curMount := range mountsStr {
		log.Println(curMount)
	}
}

func printSnapshotsOfs3Mount(mount s3wrapper.S3Mount) {
	for _, curSnapshot := range mount.Snapshots {
		log.Println(curSnapshot.Tag)
	}
}

func List(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, dbackArgs []string, spacetracker *spacetracker.SpaceTracker) {
	prefix := ``

	if len(dbackArgs) > 0 {
		prefix = strings.TrimPrefix(dbackArgs[0], `/`)
	}

	s3Mounts := getS3MountsForRestore(s3, resticw, dockerw, prefix)

	if len(s3Mounts) == 0 {
		log.Println(`Zero mounts found. Try to call "dback ls" without arguments, and then check the bucket`)
		return
	}

	if len(s3Mounts) == 1 && (s3Mounts[0].ContainerName+s3Mounts[0].Dest == dbackArgs[0]) {
		printSnapshotsOfs3Mount(s3Mounts[0])
		return
	}

	printSortedS3Mounts(s3Mounts)
}
