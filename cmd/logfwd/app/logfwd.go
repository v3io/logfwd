package app

import (
	"github.com/nuclio/logger"
	"github.com/nuclio/zap"
	"github.com/spf13/cobra"
	"github.com/v3io/go-errors"
)

type RootCommand struct {
	log     logger.Logger
	verbose bool

	cmd *cobra.Command
}

func NewCommand() *RootCommand {

	command := &RootCommand{}
	command.cmd = &cobra.Command{
		Use:   "logfwd [command]",
		Short: "Log Forwarding Service",
	}

	command.cmd.PersistentFlags().BoolVarP(&command.verbose, "verbose", "v", false, "Verbose output")
	command.cmd.AddCommand(newServerCommand(command).cmd)

	return command
}

func (c *RootCommand) Run() error {
	return c.cmd.Execute()
}

func (c *RootCommand) init() error {
	if err := c.initLogger(); err != nil {
		return errors.Wrap(err, "Unable to create logger instance")
	}

	c.log.Debug("Root Command initialized")
	return nil
}

func (c *RootCommand) initLogger() error {
	level := nucliozap.InfoLevel

	if c.verbose {
		level = nucliozap.DebugLevel
	}

	log, err := nucliozap.NewNuclioZapCmd("logfwd", level)
	if err != nil {
		return errors.Wrap(err, "Failed to create logger")
	}
	c.log = log

	return nil
}
