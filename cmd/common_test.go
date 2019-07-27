package wotwhb

import (
	"path/filepath"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type BaseSuite struct {
	workingDir string
}

var _ = Suite(&BaseSuite{})

func (s *BaseSuite) SetUpSuite(c *C) {
	s.workingDir = c.MkDir()
	ConfigDirectoryFlagValue = filepath.Join(s.workingDir, "config")
	DownloadDirectoryFlagValue = filepath.Join(s.workingDir, "downloads")
}
