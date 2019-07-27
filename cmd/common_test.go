package wotwhb

import (
	"path/filepath"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type BaseSuite struct {
	WorkingDir string
}

var _ = Suite(&BaseSuite{})

func (s *BaseSuite) SetUpSuite(c *C) {
	s.WorkingDir = c.MkDir()
	ConfigDirectoryFlagValue = filepath.Join(s.WorkingDir, "config")
	DownloadDirectoryFlagValue = filepath.Join(s.WorkingDir, "downloads")
}
