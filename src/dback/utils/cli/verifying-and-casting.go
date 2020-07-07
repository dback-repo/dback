package cli

import (
	"cli"
	"dback/utils/dockerwrapper"
	"dback/utils/resticwrapper"
	"dback/utils/s3opts"
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

const minRestoreArgs = 2
const maxRestoreArgs = 3
const maxLsArgs = 2
const minLsSnapshotsArgs = 3

func VerifyAndCast(req cli.Request) (DbackOpts, []string, resticwrapper.CreationOpts) {
	f := req.Flags

	switch req.Command {
	case `backup`:
		if len(req.Args) > 0 {
			log.Fatalln(`"dback backup" accepts no arguments`)
		}
	case `ls`:
		if len(req.Args) > 0 {
			// if len(req.Args) > maxLsArgs {
			// 	log.Fatalln(`"dback ls ` + req.Args[0] +
			// 		`" require zero, one, or two arguments. Too many arguments provided`)
			// 	if req.Args[0] == `snapshots` && len(req.Args) < minLsSnapshotsArgs {
			// 		log.Fatalln(`"dback ls snapshots" require an argument, but no arguments provided`)
			// 	}
			// }
		}
	case `restore`:
		if len(req.Args) > 0 {
			if req.Args[0] != `container` && req.Args[0] != `mount` {
				log.Fatalln(`"dback restore" accepts only subcommands "container" and "mount"`)
			}

			if len(req.Args) < minRestoreArgs {
				log.Fatalln(`"dback restore ` + req.Args[0] +
					`" require one or two arguments, but no arguments provided`)
			}

			if len(req.Args) > maxRestoreArgs {
				log.Fatalln(`"dback restore ` + req.Args[0] +
					`" require one or two arguments. Too many arguments provided`)
			}
		}
	case ``:
		//no command provided. Parse CLI is already printed an advice
		os.Exit(1)
	default:
		log.Fatalln(`Unrecognized command ` + req.Command +
			`. Run "dback --help", for list of available commands`)
	}

	dbackOpts := DbackOpts{
		Matchers:        f[`matcher`],
		ExcludePatterns: dockerwrapper.NewExcludePatterns(f[`exclude`]),
		ThreadsCount:    verifyThreads(f[`threads`][0]),
	}

	resticOpts := resticwrapper.CreationOpts{
		S3Opts: s3opts.CreationOpts{
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
			log.Fatalln(`All of dback options (s3-endpoint, s3-bucket, s3-acc-key, s3-sec-key, restic-pass)
are required, when emulation flag is not defined. Define dback options, or use emulation "-e" flag.
Run "dback ` + req.Command + ` --help", for details`)
		}
	}

	return dbackOpts, req.Args, resticOpts
}

//Positive or zero integer
func verifyThreads(str string) int {
	res, err := strconv.Atoi(str)
	if err != nil || res < 0 {
		log.Fatalln(`threads count must be zero or positive integer value`)
	}

	return res
}
