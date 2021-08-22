module github.com/0chain/zwalletcli

require (
	github.com/0chain/gosdk v1.2.82-0.20210821153536-9977518a2256
	github.com/mattn/go-colorable v0.1.4 // indirect
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
