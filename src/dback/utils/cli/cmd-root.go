package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "dback",
		Short: "Make backup of all mounts for all of containers",
		Long: `Dback is application for observe docker containers, make bulk incremental backups 
of their mounts (folders and volumes), and pass backups to S3 bucket.
Dback runs restic under the hood.

Main options:
- Filter containers bulk backup by container properties.
- Exclude containers and mounts by name/path (regex).

Also you can restore backups to exist containers.`,
	}
	emulate := false
	c.PersistentFlags().BoolVarP(&emulate, "emulate", "e", false,
		"emulate an action (backup or restore), and show list of items will be affected")

	return &c
}
