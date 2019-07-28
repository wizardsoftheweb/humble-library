package wotwhb

import (
	"github.com/spf13/cobra"
)

func init() {
	PackageCmd.AddCommand(UpdateCmd)
	UpdateCmd.AddCommand(UpdateKeyListCmd)
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

var UpdateKeyListCmd = &cobra.Command{
	Use:   "keys",
	Short: "Update the full list of purchases (bundles, store purchases, etc.)",
	Run:   UpdateKeyListCmdRun,
}

func UpdateKeyListCmdRun(cmd *cobra.Command, args []string) {
	updateKeyList(cmd)
}

var UpdateOrderListCmd = &cobra.Command{
	Use:   "orders",
	Short: "Update information about all the available orders (bundle contents, file type variants, links, etc.)",
	Run:   UpdateOrderListCmdRun,
}

func UpdateOrderListCmdRun(cmd *cobra.Command, args []string) {
	keys := loadSavedKeyList()
	updateAllOrders(cmd, keys)
}
