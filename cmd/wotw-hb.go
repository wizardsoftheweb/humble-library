package wotwhb

import (
	"github.com/spf13/cobra"
)

var PackageVersion = "0.0.0"
var VerbosityFlagValue int
var ConfigDirectoryFlagValue string
var DownloadDirectoryFlagValue string

func init() {
	PackageCmd.PersistentFlags().CountVarP(
		&VerbosityFlagValue,
		"verbose",
		"v",
		"Increases application verbosity",
	)
	PackageCmd.PersistentFlags().StringVarP(
		&ConfigDirectoryFlagValue,
		"config-directory",
		"c",
		configDirectory,
		"Location to store configuration files",
	)
	PackageCmd.PersistentFlags().StringVarP(
		&DownloadDirectoryFlagValue,
		"download-directory",
		"d",
		downloadDirectory,
		"Location to store downloads",
	)
}

func Execute() error {
	return PackageCmd.Execute()
}

var PackageCmd = &cobra.Command{
	Use:              "wotw-hb",
	Version:          PackageVersion,
	Short:            "WIP Humble Bundle CLI",
	PersistentPreRun: PackageCmdPersistentPreRun,
	Run:              HelpOnly,
}

func PackageCmdPersistentPreRun(cmd *cobra.Command, args []string) {
	BootstrapLogger(VerbosityFlagValue)
	BootstrapConfig(ConfigDirectoryFlagValue, DownloadDirectoryFlagValue)
}

func HelpOnly(cmd *cobra.Command, args []string) {
	RunUi()
	// _ = cmd.Help()
}
