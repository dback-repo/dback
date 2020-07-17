package cli

import (
	"cli"

	"github.com/spf13/cobra"
)

func NewInspectCommand(reqest *cli.Request) *cobra.Command {
	c := cobra.Command{
		Use:   "inspect",
		Short: "Prints 'docker inspect' converted to XML",
		Long: `Prints 'docker inspect' converted to XML. Use it to write backup matchers.
  dback inspect <container>`,
		Run: func(cmd *cobra.Command, args []string) {
			reqest.Command = cmd.Use
			reqest.Args = args
		},
	}

	return &c
}
