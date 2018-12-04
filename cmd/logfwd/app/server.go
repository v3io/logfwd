package app

import (
	"github.com/v3io/logfwd/pkg/rules"
	rest "github.com/v3io/logfwd/pkg/server"
	"github.com/spf13/cobra"
	"github.com/v3io/go-errors"
)

type serverCommand struct {
	root *RootCommand

	cmd           *cobra.Command
	listenAddress string
}

func newServerCommand(root *RootCommand) *serverCommand {
	command := &serverCommand{
		root: root,
	}
	cmd := &cobra.Command{
		Use:   "server rules-config-path [options]",
		Short: "Run forwarding server with rules configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Missing rules config path")
			}

			if err := root.init(); err != nil {
				return errors.Wrap(err, "Failed to initialize root command")
			}
			serverLog := root.log.GetChild("cmd.server")
			serverLog.DebugWith("Creating rules configuration", "rules-file", args[0])
			rulesConfig, configErr := rules.NewRuleConfig(root.log, args[0])
			if configErr != nil {
				return errors.Wrap(configErr, "Unable to read rules-file")
			}
			serverLog.InfoWith("Creating REST server", "listenAddress", command.listenAddress)
			server, serverErr := rest.NewServer(root.log, command.listenAddress, rulesConfig)
			if serverErr != nil {
				return errors.Wrap(serverErr, "Unable to create server")
			}

			serverLog.Info("Running REST server")
			return server.Run()
		},
	}

	cmd.Flags().StringVarP(&command.listenAddress,
		"listen-address",
		"s",
		"0.0.0.0:8080",
		"Listen address for server")

	command.cmd = cmd
	return command
}
