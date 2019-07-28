package wotwhb

import (
	. "gopkg.in/check.v1"
)

type LocalDataSuite struct {
	BaseSuite
}

var _ = Suite(&LocalDataSuite{})

func (s *LocalDataSuite) Printf(format string, args ...interface{}) {}

func (s *LocalDataSuite) TestParseRawKeyList(c *C) {
	testData := []byte(`[{"gamekey":"one"},{"gamekey":"two"},{"gamekey":"three"}]`)
	keys := parseRawKeyList(testData)
	c.Assert(keys, DeepEquals, savedKeyListTest)
}

func (s *LocalDataSuite) TestLoadSavedKeyList(c *C) {
	results := loadSavedKeyList()
	c.Assert(results, DeepEquals, savedKeyListTest)
}

func (s *LocalDataSuite) TestQueryAllOrders(c *C) {
	c.Assert(
		func() {
			queryAllOrders(s, savedKeyListTest)
		},
		PanicMatches,
		".*JSON.*",
	)
	result := queryAllOrders(s, []string{})
	c.Assert(len(result), Equals, 0)

}

func (s *LocalDataSuite) TestLoadAllOrdersAsStruct(c *C) {
	result := loadAllOrdersAsStruct()
	c.Assert(len(result), Equals, 0)
}
