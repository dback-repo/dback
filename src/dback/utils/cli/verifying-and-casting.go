package cli

import (
	"cli"
	"dback/utils/dockerwrapper"
)

func VerifyAndCast(req cli.Request) (dockerwrapper.EmulateFlag, []dockerwrapper.ExcludePattern) {
	f := req.Flags
	excludePatterns := dockerwrapper.NewExcludePatterns(f[`exclude`])
	isEmulation := dockerwrapper.NewEmulateFlag(f[`emulate`][0])

	return isEmulation, excludePatterns
}
