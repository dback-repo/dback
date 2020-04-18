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
	fx := func(cmd *cobra.Command, args []string) {
		reqest.Command = `backup`
		reqest.Flags[`emulate`] = cmd.Flag(`emulate`).Value.String()
		reqest.Args = args
	}

	res := cliParser{}
	res.rootCmd = cobra.Command{
		Use:   "cobraex",
		Short: "Make backup of all mounts for all of containers",
		Long: `Dback is application for observe docker containers, do incremental backups 
of their mounts (bind and volumes), and pass it to S3 bucket. 
You able to exclude an extra data from observed backup list.
Also you can restore backups to exist mounts.
Dback runs restic under the hood.`,
	}

	res.rootCmd.AddCommand(&cobra.Command{
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
		Run: fx,
	})

	emulate := false
	res.rootCmd.PersistentFlags().BoolVarP(&emulate, "emulate", "e", false,
		"emulate an action, and show list of items will be affected")

	return &res
}

func (t *cliParser) Parse() Request {
	check(t.rootCmd.Execute())
	return Request{}
}

func ParseCLI() Request {
	res := NewRequest()
	newCliParser(&res).Parse()

	return res
}
