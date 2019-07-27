package wotwhb

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var PackageVersion = "0.0.0"
var VerbosityFlagValue int

func init() {
	PackageCmd.PersistentFlags().CountVarP(
		&VerbosityFlagValue,
		"verbose",
		"v",
		"Increases application verbosity",
	)
}

func Execute() {
	if err := PackageCmd.Execute(); nil != err {
		logrus.Fatal(err)
	}
}

var PackageCmd = &cobra.Command{
	Use:              "wotw-hb",
	Version:          PackageVersion,
	Short:            "WIP Humble Bundle CLI",
	PersistentPreRun: PackageCmdPersistentPreRun,
	Run:              PackageCmdRun,
}

func PackageCmdPersistentPreRun(cmd *cobra.Command, args []string) {
	BootstrapLogger(VerbosityFlagValue)
}

func PackageCmdRun(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
