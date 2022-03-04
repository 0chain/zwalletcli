package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/icza/bitio"
	"github.com/spf13/cobra"
)

var VersionStr string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version information",
	Long:  `Prints version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version info:")
		fmt.Println("\tzwallet...: ", VersionStr)
		fmt.Println("\tgosdk.....: ", getVersion("github.com/0chain/gosdk"))
		return
	},
}

func getVersion(path string) string {
	_ = bitio.NewReader
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Failed to read build info")
		return ""
	}

	for _, dep := range bi.Deps {
		if dep.Path == path {
			if dep.Replace != nil && dep.Replace.Version != "" {
				return dep.Replace.Version
			}

			return dep.Version
		}
	}

	return ""
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
