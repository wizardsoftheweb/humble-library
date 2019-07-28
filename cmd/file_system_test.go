package wotwhb

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type FileSystemSuite struct {
	BaseSuite
}

var _ = Suite(&FileSystemSuite{})

func (s *FileSystemSuite) TestEnsureDirectoryExists(c *C) {
	directoryName := filepath.Join(s.WorkingDir, "test")
	_, err := os.Stat(directoryName)
	c.Assert(err, NotNil)
	ensureDirectoryExists(directoryName)
	_, err = os.Stat(directoryName)
	c.Assert(err, IsNil)
}

func (s *FileSystemSuite) TestWriteJson(c *C) {
	testContents := []int{1, 2, 3}
	fileName := filepath.Join(s.WorkingDir, "test.json")
	writeJsonToFile(testContents, fileName)
	contents, _ := ioutil.ReadFile(fileName)
	var keys []int
	err := json.Unmarshal(contents, &keys)
	c.Assert(err, IsNil)
	c.Assert(testContents, DeepEquals, keys)
}
