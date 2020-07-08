package cli

import (
	"cli"

	"github.com/spf13/cobra"
)

func NewListCommand(reqest *cli.Request) *cobra.Command {
	c := cobra.Command{
		Use:   "ls",
		Short: "List present mounts, and snapshots in the bucket",
		Long: `List present mounts, and snapshots in the bucket.
With no arguments, "ls" prints all the mounts present in the bucket.
  dback ls <...flags>

If you provide a mount prefix as argument - "ls" show mounts starts with prefix.
For example - mounts of some container.
  dback ls /app-container-name <...flags>

If you provided full name of a mount - "ls" show snapshots of the mount.
You are able to restore snapshot you need with "dback restore" command.
  dback ls /app-container-name/var/lib/mysql <...flags>
`,
		Run: func(cmd *cobra.Command, args []string) {
			reqest.Command = cmd.Use
			reqest.Flags[`threads`] = []string{cmd.Flag(`threads`).Value.String()}
			reqest.Flags[`s3-endpoint`] = []string{cmd.Flag(`s3-endpoint`).Value.String()}
			reqest.Flags[`s3-bucket`] = []string{cmd.Flag(`s3-bucket`).Value.String()}
			reqest.Flags[`s3-acc-key`] = []string{cmd.Flag(`s3-acc-key`).Value.String()}
			reqest.Flags[`s3-sec-key`] = []string{cmd.Flag(`s3-sec-key`).Value.String()}
			reqest.Flags[`restic-pass`] = []string{cmd.Flag(`restic-pass`).Value.String()}

			reqest.Args = args
		},
	}

	c.PersistentFlags().StringP(`threads`, `t`, `0`, `run mounts backup concurrently. 0 - create a thread for each mount`)
	c.PersistentFlags().String(`s3-endpoint`, ``, `with protocol and port "http://192.168.0.3:1337"`)
	c.PersistentFlags().StringP(`s3-bucket`, `b`, ``, `name of bucket at the s3 endpoint`)
	c.PersistentFlags().StringP(`s3-acc-key`, `a`, ``, `s3 access key`)
	c.PersistentFlags().StringP(`s3-sec-key`, `s`, ``, `s3 secret key`)
	c.PersistentFlags().StringP(`restic-pass`, `p`, ``, `restic password`)

	return &c
}
