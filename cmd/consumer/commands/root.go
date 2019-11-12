package commands

import (
	"github.com/fusion-app/fusion-app/cmd/consumer/commands/subscriber"
	"github.com/spf13/cobra"
)

const (
	// CLIName is the name of this CLI
	CLIName = "fusionctl"
)

func NewCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   CLIName,
		Short: "fusionctl is the command line interface to fusion-app",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(subscriber.NewSubscribeCommand())

	return command
}
