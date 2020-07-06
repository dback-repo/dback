package cli

import (
	"cli"

	"github.com/spf13/cobra"
)

func NewRestoreCommand(reqest *cli.Request) *cobra.Command {
	c := cobra.Command{
		Use:   "restore",
		Short: "Restore backups from s3 to exist containers or mounts",
		Long: `Restore backups from s3 to exist containers or mounts.

Container restore 
	Restore a snapshot to local mounts of container, wich has the same name
	dback restore container app-container-name <...flags>

Container restore, with different name
	Restore a snapshot to local mounts of container, wich has a different name
	dback restore container app-container-name app-new-local-container-name <...flags>

Mount restore
	Restore a snapshot to local mount with the same container name
	dback restore mount app-container-name/var/lib/mysql <...flags>

Mount restore, with different container name
	Restore a snapshot to local mount with the same container name
	dback restore mount app-container-name/var/lib/mysql app-new-local-container-name/var/lib/mysql <...flags>

Bulk restore
	Find all backups in s3 bucket, then restore all mounts exist at the host
	dback restore <...flags>
`,
		Run: func(cmd *cobra.Command, args []string) {
			reqest.Command = cmd.Use
			reqest.Flags[`emulate`] = []string{cmd.Flag(`emulate`).Value.String()}
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
