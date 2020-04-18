package cli

// import (
// 	"fmt"

// 	"github.com/spf13/cobra"
// )

// // listCmd represents the list command
// var listCmd = &cobra.Command{
// 	Hidden: true,
// 	Use:    "ls",
// 	Short:  "Make list, and pass it to S3 bucket",
// 	Long:   `Longer description: Make list, and pass it to S3 bucket`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("list called")
// 		fmt.Println(args)
// 		fmt.Println(cmd.Flag(`emulate`).Value.String())
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(listCmd)
// }

// command string
// flags map(string)string
// params []string

type clirequest struct {
}
