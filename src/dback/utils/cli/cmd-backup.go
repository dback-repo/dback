package cli

import (
	"cli"

	"github.com/spf13/cobra"
)

func NewBackupCommand(reqest *cli.Request) *cobra.Command {
	c := cobra.Command{
		Use:   "backup",
		Short: "Observe containers, make mounts backup, then pass backups to S3 bucket",
		Long: `Observe containers, make mounts backup, then pass backups to S3 bucket
  Create snapshot of mounts matched all the points:
  default points:
  - HostConfig.AutoRemove:      false
  - HostConfig.RestartPolicy:   always

Options:
  --exclude            Exclude volume pattern
    mounts are named as: [ContainerName]/[PathInContainer]
    For example, mount in "mysql" container: mysql/var/mysql/data
    Pattern is a regular expression. For example, "^/(drone.*|dback-test-1.5.*)$"
    ignore all mounts starts with "/drone", or "/dback-test-1.5"`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			reqest.Command = cmd.Use
			reqest.Flags[`emulate`] = []string{cmd.Flag(`emulate`).Value.String()}

			reqest.Flags[`exclude`], err = cmd.PersistentFlags().GetStringSlice(`exclude`)
			check(err)

			reqest.Flags[`matcher`], err = cmd.PersistentFlags().GetStringSlice(`matcher`)
			check(err)

			reqest.Flags[`threads`] = []string{cmd.Flag(`threads`).Value.String()}
			reqest.Flags[`s3-endpoint`] = []string{cmd.Flag(`s3-endpoint`).Value.String()}
			reqest.Flags[`s3-bucket`] = []string{cmd.Flag(`s3-bucket`).Value.String()}
			reqest.Flags[`s3-acc-key`] = []string{cmd.Flag(`s3-acc-key`).Value.String()}
			reqest.Flags[`s3-sec-key`] = []string{cmd.Flag(`s3-sec-key`).Value.String()}
			reqest.Flags[`s3-region`] = []string{cmd.Flag(`s3-region`).Value.String()}
			reqest.Flags[`restic-pass`] = []string{cmd.Flag(`restic-pass`).Value.String()}

			reqest.Args = args
		},
	}
	c.PersistentFlags().StringSliceP(`matcher`, `m`,
		[]string{`"RestartPolicy":{"Name":"always"`, `"AutoRemove":false`, `"Running":true`},
		`backup containers with all defined preferences`)
	c.PersistentFlags().StringSliceP(`exclude`, `x`, []string{}, `exclude mounts by RegEx pattern`)
	c.PersistentFlags().StringP(`threads`, `t`, `0`, `run mounts backup concurrently. 0 - create a thread for each mount`)
	c.PersistentFlags().String(`s3-endpoint`, ``, `with protocol and port "http://192.168.0.3:1337"`)
	c.PersistentFlags().StringP(`s3-bucket`, `b`, ``, `name of bucket at the s3 endpoint`)
	c.PersistentFlags().StringP(`s3-acc-key`, `a`, ``, `s3 access key`)
	c.PersistentFlags().StringP(`s3-sec-key`, `s`, ``, `s3 secret key`)
	c.PersistentFlags().String(`s3-region`, ``, `s3 region, for example eu-central-1`)
	c.PersistentFlags().StringP(`restic-pass`, `p`, ``, `restic password`)

	return &c
}
