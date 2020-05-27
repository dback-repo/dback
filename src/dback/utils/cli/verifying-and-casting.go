package cli

import (
	"cli"
	"dback/utils/dockerwrapper"
	"log"
	"strconv"
)

func VerifyAndCast(req cli.Request) (dockerwrapper.EmulateFlag, []dockerwrapper.ExcludePattern,
	int, string, string, string, string, string) {
	f := req.Flags
	excludePatterns := dockerwrapper.NewExcludePatterns(f[`exclude`])
	isEmulation := dockerwrapper.NewEmulateFlag(f[`emulate`][0])

	threads := verifyThreads(f[`threads`][0])
	s3Endpoint := f[`s3-endpoint`][0]
	s3Bucket := f[`s3-bucket`][0]
	s3AccKey := f[`s3-acc-key`][0]
	s3SecKey := f[`s3-sec-key`][0]
	resticPass := f[`restic-pass`][0]

	return isEmulation, excludePatterns, threads, s3Endpoint, s3Bucket, s3AccKey, s3SecKey, resticPass
}

//Positive integer
func verifyThreads(str string) int {
	res, err := strconv.Atoi(str)
	if err != nil || res < 0 {
		log.Fatalln(`threads count must be zero or positive integer value`)
	}

	return res
}
