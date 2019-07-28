package wotwhb

import (
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update local databases",
}

func UpdateCmdRun(cmd *cobra.Command, args []string) {
	cmd.Println("rad")
}
