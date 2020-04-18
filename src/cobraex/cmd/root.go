package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cobraex",
	Short: "Make backup of all mounts for all of containers",
	Long: `Dback is application for observe docker containers, do incremental backups 
of their mounts (bind and volumes), and pass it to S3 bucket. 
You able to exclude an extra data from observed backup list.
Also you can restore backups to exist mounts.
Dback runs restic under the hood.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var emulate bool

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobraex.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("emulate", "e", false, "emulate write actions, and show list of items that will be affected")
	//cobra.MarkFlagRequired()

	rootCmd.PersistentFlags().BoolVarP(&emulate, "emulate", "e", false, "emulate an action, and show list of items will be affected")
	//rootCmd.MarkPersistentFlagRequired(`emulate`)
	//fmt.Println(emulate)
	//rootCmd.ParseFlags()

	//rootCmd.Flags()
	//rootCmd.Flags().

	// rootCmd.AddCommand(backupCmd)
	// rootCmd.AddCommand(restoreCmd)
}
