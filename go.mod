module github.com/0chain/zwalletcmd

require (
	github.com/0chain/gosdk v0.0.0-00010101000000-000000000000
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	gopkg.in/cheggaaa/pb.v1 v1.0.28
)

replace github.com/0chain/zwalletcmd => ../zwalletcmd

replace github.com/0chain/gosdk => ../gosdk

go 1.12
