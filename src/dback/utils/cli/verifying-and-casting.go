package cli

import (
	"cli"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3wrapper"
	"log"
	"os"
	"strconv"
)

type DbackOpts struct {
	IsEmulation     dockerwrapper.EmulateFlag
	Matchers        []string
	ExcludePatterns []dockerwrapper.ExcludePattern
	ThreadsCount    int
}

func VerifyAndCast(req cli.Request) (DbackOpts, resticwrapper.CreationOpts) {
	f := req.Flags

	switch req.Command {
	case `backup`:
		if len(req.Args) > 0 {
			log.Fatalln(`"dback backup" accepts no arguments`)
		}
	case ``:
		//no command provided. Parse CLI is already printed an advice
		os.Exit(1)
	default:
		log.Fatalln(`Unrecognized command ` + req.Command +
			`. Run "dback --help", for list of available commands`)
	}

	dbackOpts := DbackOpts{
		IsEmulation:     dockerwrapper.NewEmulateFlag(f[`emulate`][0]),
		Matchers:        f[`matcher`],
		ExcludePatterns: dockerwrapper.NewExcludePatterns(f[`exclude`]),
		ThreadsCount:    verifyThreads(f[`threads`][0]),
	}

	resticOpts := resticwrapper.CreationOpts{
		S3Opts: s3wrapper.CreationOpts{
			S3Endpoint: f[`s3-endpoint`][0],
			S3Bucket:   f[`s3-bucket`][0],
			AccKey:     f[`s3-acc-key`][0],
			SecKey:     f[`s3-sec-key`][0],
		},
		ResticPass: f[`restic-pass`][0],
	}

	if !dbackOpts.IsEmulation {
		if resticOpts.S3Opts.S3Endpoint == `` || resticOpts.S3Opts.S3Bucket == `` ||
			resticOpts.S3Opts.AccKey == `` || resticOpts.S3Opts.SecKey == `` ||
			resticOpts.ResticPass == `` {
			log.Fatalln(`All of restic options (s3-endpoint, s3-bucket, s3-acc-key, s3-sec-key, restic-pass)
are required, when emulation flag is not defined. Define restic options, or use emulation "-e" flag.
Run "dback backup --help", for details`)
		}
	}

	return dbackOpts, resticOpts
}

//Positive or zero integer
func verifyThreads(str string) int {
	res, err := strconv.Atoi(str)
	if err != nil || res < 0 {
		log.Fatalln(`threads count must be zero or positive integer value`)
	}

	return res
}
