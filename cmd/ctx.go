package cmd

import "github.com/spf13/cobra"

var (
	withoutZCNCoreCmds = make(map[*cobra.Command]bool)
	withoutWalletCmds  = make(map[*cobra.Command]bool)
)

// WithoutZCNCore zcncore package is unnecessary for this command. it will be asked to initialize zcncore via zcncore.Init
func WithoutZCNCore(c *cobra.Command) *cobra.Command {
	withoutZCNCoreCmds[c] = true
	return c
}

// WithoutWallet wallet information is unnecessary for this command. ~/.zcn/wallet.json will not be checked, and wallet will not be asked to register on blockchain
func WithoutWallet(c *cobra.Command) *cobra.Command {
	withoutWalletCmds[c] = true
	return c
}
