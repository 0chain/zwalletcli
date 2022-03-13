module github.com/0chain/zwalletcli

require (
	github.com/0chain/gosdk v1.7.7-0.20220313102505-53d2cafe34a9
	github.com/olekukonko/tablewriter v0.0.5
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.9.0
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

go 1.16

// temporary, for development
//replace github.com/0chain/gosdk => ../gosdk
