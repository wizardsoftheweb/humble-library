package wotwhb

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wotw-hb",
	Short: "WIP Humble Bundle CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("rad")
	},
}

func Execute() {
	if err := rootCmd.Execute(); nil != err {
		logrus.Fatal(err)
	}
}
