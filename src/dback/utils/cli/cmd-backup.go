package cli

import (
	"cli"

	"github.com/spf13/cobra"
)

func NewBackupCommand(reqest *cli.Request) *cobra.Command {
	c := cobra.Command{
		Use:   "backup",
		Short: "Make backup, and pass it to S3 bucket",
		Long: `Make backup, and pass it to S3 bucket
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
			reqest.Command = cmd.Use
			reqest.Flags[`emulate`] = []string{cmd.Flag(`emulate`).Value.String()}
			var err error
			reqest.Flags[`exclude`], err = cmd.PersistentFlags().GetStringSlice(`exclude`)
			check(err)
			reqest.Args = args
		},
	}
	c.PersistentFlags().StringSliceP("exclude", "x", []string{}, "exclude containers by name, matched with RegEx")

	return &c
}
