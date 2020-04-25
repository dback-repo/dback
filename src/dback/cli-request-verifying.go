package main

import (
	"dback/utils/cli"
	"dback/utils/dockerwrapper"
)

func verifyCliReq(req cli.Request) (dockerwrapper.EmulateFlag, []dockerwrapper.ExcludePattern) {
	f := req.Flags
	excludePatterns := dockerwrapper.NewExcludePatterns(f[`exclude`])
	isEmulation := dockerwrapper.NewEmulateFlag(f[`emulate`][0])

	return isEmulation, excludePatterns
}
