package cmd

import "github.com/spf13/cobra"

var (
	withoutZCNCommands    = make(map[*cobra.Command]bool)
	withoutWalletCommands = make(map[*cobra.Command]bool)
)

func WithoutZCN(c *cobra.Command) *cobra.Command {
	withoutZCNCommands[c] = true
	return c
}

func WithoutWallet(c *cobra.Command) *cobra.Command {
	withoutWalletCommands[c] = true
	return c
}
