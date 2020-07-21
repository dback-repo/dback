package cli

import (
	"cli"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ParseCLI() cli.Request {
	cliRequest := cli.NewRequest()
	cli.ParseCLI(
		NewRootCommand(),
		NewBackupCommand(&cliRequest),
		NewListCommand(&cliRequest),
		NewRestoreCommand(&cliRequest),
		NewInspectCommand(&cliRequest))

	return cliRequest
}
