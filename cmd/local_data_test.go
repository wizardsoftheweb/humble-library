package wotwhb

import (
	"path/filepath"

	. "gopkg.in/check.v1"
)

type LocalDataSuite struct {
	BaseSuite
}

var savedKeyListTest = []string{"one", "two", "three"}

var _ = Suite(&LocalDataSuite{})

func (s *LocalDataSuite) SetUpTest(c *C) {
	writeJsonToFile(savedKeyListTest, filepath.Join(ConfigDirectoryFlagValue, orderKeyListFileBasename))
}

func (s *LocalDataSuite) TestParseRawKeyList(c *C) {
	testData := []byte(`[{"gamekey":"one"},{"gamekey":"two"},{"gamekey":"three"}]`)
	keys := parseRawKeyList(testData)
	c.Assert(keys, DeepEquals, savedKeyListTest)
}

func (s *LocalDataSuite) TestLoadSavedKeyList(c *C) {
	results := loadSavedKeyList()
	c.Assert(results, DeepEquals, savedKeyListTest)
}
