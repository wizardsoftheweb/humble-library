package wotwhb

import (
	"github.com/spf13/cobra"
)

func init() {
	PackageCmd.AddCommand(UpdateCmd)
	UpdateCmd.AddCommand(UpdateOrderListCmd)
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update local databases",
	Run:   UpdateCmdRun,
}

func UpdateCmdRun(cmd *cobra.Command, args []string) {
	cmd.Println("rad")
}

var UpdateOrderListCmd = &cobra.Command{
	Use:   "orders",
	Short: "Update the full list of purchases (bundles, store purchases, etc.)",
	Run:   UpdateOrderListCmdRun,
}

func UpdateOrderListCmdRun(cmd *cobra.Command, args []string) {
	updateOrderList(cmd)
}
