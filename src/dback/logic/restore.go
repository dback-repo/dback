package logic

import (
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"log"
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

	log.Println(`localMounts:`, localMounts)
	log.Println(`restoreParams:`, restoreParams)

	return res
}

func Restore(s3 *s3wrapper.S3Wrapper, resticw *resticwrapper.ResticWrapper,
	dockerw *dockerwrapper.DockerWrapper, restoreParams []string) {
	s3MountsForRestore := getS3MountsForRestore(s3, resticw, dockerw, restoreParams)

	log.Println(s3MountsForRestore)
}
