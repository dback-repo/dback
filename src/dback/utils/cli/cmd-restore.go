package cli

import (
	"github.com/spf13/cobra"
)

func NewRestoreCommand(reqest *Request) *cobra.Command {
	return &cobra.Command{
		Use:   "restore",
		Short: "restore backups from s3 to exist mounts",
		Long:  `TOODO: longer description of restore`,
		Run: func(cmd *cobra.Command, args []string) {
			reqest.Command = cmd.Use
			reqest.Flags[`emulate`] = []string{cmd.Flag(`emulate`).Value.String()}
		},
	}
}
