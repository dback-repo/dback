package cli

import (
	"github.com/spf13/cobra"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Request struct {
	Command string
	Flags   map[string][]string
	Args    []string
}

func NewRequest() Request {
	return Request{``, make(map[string][]string), []string{}}
}

type cliParser struct {
	rootCmd cobra.Command
}

func NewCliParser(root *cobra.Command, commands ...*cobra.Command) *cliParser {
	res := cliParser{}
	res.rootCmd = *root

	for _, curCommand := range commands {
		res.rootCmd.AddCommand(curCommand)
	}

	return &res
}

func (t *cliParser) Parse() {
	check(t.rootCmd.Execute())
}

func ParseCLI(root *cobra.Command, commands ...*cobra.Command) Request {
	res := NewRequest()
	NewCliParser(root, commands...).Parse()

	return res
}
