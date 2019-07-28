package wotwhb

import (
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
