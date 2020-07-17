package logic

import (
	"dback/utils/dockerwrapper"
	"fmt"
)

func Inspect(dockerw *dockerwrapper.DockerWrapper, dbackArgs []string) {
	fmt.Println(dockerw.GetDockerInspectXMLByContainerName(containerNameLeadingSlash(dbackArgs[0])))
}
