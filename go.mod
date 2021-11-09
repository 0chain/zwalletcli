module github.com/0chain/zwalletcli

require (
	github.com/0chain/gosdk v0.0.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

go 1.13

// temporary, for development
// replace github.com/0chain/gosdk => ../gosdk
