package wotwhb

import (
	. "gopkg.in/check.v1"
)

type MungeSuite struct {
	BaseSuite
}

var _ = Suite(&MungeSuite{})

func (s *MungeSuite) TestTraverseAllOrders(c *C) {
	var orders []HbOrder
	orders = []HbOrder{
		{
			AmountSpent: float64(0),
			SubProducts: []HbSubProduct{
				{
					HumanName: "qqq",
					Downloads: []HbSubProductDownload{
						{
							DownloadStruct: []HbSubProductDownloadDownloadStruct{
								{
									Arch: "32",
								},
							},
						},
					},
				},
			},
		},
	}
	results := traverseAllOrders(orders)
	c.Assert(results.TotalFileSize, Equals, int64(0))
}
