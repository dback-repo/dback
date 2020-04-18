/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	Long:  `Longer description: Make backup, and pass it to S3 bucket`,

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
