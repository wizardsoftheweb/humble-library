package wotwhb

import (
	"path/filepath"

	. "gopkg.in/check.v1"
)

type LocalDataSuite struct {
	BaseSuite
}

var savedOrderListTest = []string{"one", "two", "three"}

var _ = Suite(&LocalDataSuite{})

func (s *LocalDataSuite) SetUpTest(c *C) {
	writeJsonToFile(savedOrderListTest, filepath.Join(ConfigDirectoryFlagValue, orderKeyListFileBasename))
}

func (s *LocalDataSuite) TestParseRawOrderList(c *C) {
	testData := []byte(`[{"gamekey":"one"},{"gamekey":"two"},{"gamekey":"three"}]`)
	keys := parseRawOrderList(testData)
	c.Assert(keys, DeepEquals, savedOrderListTest)
}

func (s *LocalDataSuite) TestLoadSavedOrderList(c *C) {
	results := loadSavedOrderList()
	c.Assert(results, DeepEquals, savedOrderListTest)
}
