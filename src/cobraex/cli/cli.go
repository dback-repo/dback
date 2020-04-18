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
	Flags   map[string]string
	Args    []string
}

func NewRequest() Request {
	return Request{``, make(map[string]string), []string{}}
}

type cliParser struct {
	rootCmd cobra.Command
}

func newCliParser(reqest *Request) *cliParser {
	res := cliParser{}
	res.rootCmd = *NewRootCommand()
	res.rootCmd.AddCommand(NewBackupCommand(reqest), NewRestoreCommand(reqest))

	return &res
}

func (t *cliParser) Parse() {
	check(t.rootCmd.Execute())
}

func ParseCLI() Request {
	res := NewRequest()
	newCliParser(&res).Parse()

	return res
}
