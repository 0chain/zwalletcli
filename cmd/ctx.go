package cmd

import "github.com/spf13/cobra"

var (
	withoutZCNCmds    = make(map[*cobra.Command]bool)
	withoutWalletCmds = make(map[*cobra.Command]bool)
)

func WithoutZCN(c *cobra.Command) *cobra.Command {
	withoutZCNCmds[c] = true
	return c
}

func WithoutWallet(c *cobra.Command) *cobra.Command {
	withoutWalletCmds[c] = true
	return c
}
