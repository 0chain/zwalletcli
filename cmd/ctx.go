package cmd

import "github.com/spf13/cobra"

var offlineCommands = make(map[*cobra.Command]bool)

func EnableOffline(c *cobra.Command) *cobra.Command {

	offlineCommands[c] = true
	return c
}
