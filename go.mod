module github.com/0chain/zwalletcli

require (
	github.com/0chain/gosdk v1.4.1-0.20220117131359-7dd554e5a0f5
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.9.0
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

go 1.16

// temporary, for development
//replace github.com/0chain/gosdk => ../gosdk
