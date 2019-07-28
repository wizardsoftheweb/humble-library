package wotwhb

import (
	"github.com/spf13/cobra"
)

func init() {
	PackageCmd.AddCommand(UpdateCmd)
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update local databases",
	Run:   UpdateCmdRun,
}

func UpdateCmdRun(cmd *cobra.Command, args []string) {
	cmd.Println("rad")
}
