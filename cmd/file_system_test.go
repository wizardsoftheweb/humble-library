package wotwhb

import (
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type FileSystemSuite struct {
	BaseSuite
	tmp string
}

var _ = Suite(&FileSystemSuite{})

func (s *FileSystemSuite) SetUpTest(c *C) {
	s.tmp = c.MkDir()
}

func (s *FileSystemSuite) TestEnsureDirectoryExists(c *C) {
	directoryName := filepath.Join(s.tmp, "test")
	_, err := os.Stat(directoryName)
	c.Assert(err, NotNil)
	ensureDirectoryExists(directoryName)
	_, err = os.Stat(directoryName)
	c.Assert(err, IsNil)
}
