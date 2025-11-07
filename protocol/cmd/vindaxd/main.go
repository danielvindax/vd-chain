package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/danielvindax/vd-chain/protocol/app"
	"github.com/danielvindax/vd-chain/protocol/app/config"
	"github.com/danielvindax/vd-chain/protocol/app/constants"
	"github.com/danielvindax/vd-chain/protocol/cmd/vindaxd/cmd"
)

func main() {
	config.SetupConfig()

	option := cmd.GetOptionWithCustomStartCmd()
	rootCmd := cmd.NewRootCmd(option, app.DefaultNodeHome)

	cmd.AddTendermintSubcommands(rootCmd)
	cmd.AddInitCmdPostRunE(rootCmd)

	if err := svrcmd.Execute(rootCmd, constants.AppDaemonName, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
