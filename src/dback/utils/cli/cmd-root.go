package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "dback",
		Short: "Make backup of all mounts for all of containers",
		Long: `Dback is application for observe docker containers, do incremental backups 
of their mounts (bind and volumes), and pass it to S3 bucket. 
You able to exclude an extra data from observed backup list.
Also you can restore backups to exist mounts.
Dback runs restic under the hood.`,
	}
	emulate := false
	c.PersistentFlags().BoolVarP(&emulate, "emulate", "e", false,
		"emulate an action, and show list of items will be affected")

	return &c
}