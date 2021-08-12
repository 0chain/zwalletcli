module github.com/0chain/zwalletcli

require (
	github.com/0chain/gosdk v0.0.0
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/cobra v1.1.1
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

go 1.13

// temporary, for development
replace github.com/0chain/gosdk => ../gosdk
