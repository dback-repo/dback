package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var exclude string

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
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
		fmt.Println("backup called")
		fmt.Println(args)
		//fmt.Println(`emulate`, cmd.Flag(`emulate`).Value.String())
		//fmt.Println(`exclude`, cmd.Flag(`exclude`).Value.String())
		fmt.Println(`emulate`, emulate)
		fmt.Println(`exclude`, exclude)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.PersistentFlags().StringVarP(&exclude, "exclude", "x", ``, "Exclude containers by name, matched with RegEx")
}
