package wotwhb

import (
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type BaseSuite struct {
	WorkingDir string
	Command    *cobra.Command
	Args       []string
}

var savedKeyListTest = []string{"one", "two", "three"}
var _ = Suite(&BaseSuite{})

func (s *BaseSuite) SetUpSuite(c *C) {
	s.WorkingDir = c.MkDir()
	ConfigDirectoryFlagValue = filepath.Join(s.WorkingDir, "config")
	DownloadDirectoryFlagValue = filepath.Join(s.WorkingDir, "downloads")
	fatalHandler = func(args ...interface{}) { panic(args[0]) }
	BootstrapConfig(ConfigDirectoryFlagValue, DownloadDirectoryFlagValue)
	writeJsonToFile(savedKeyListTest, filepath.Join(ConfigDirectoryFlagValue, orderKeyListFileBasename))
	s.Command = &cobra.Command{}
	s.Args = []string{}
}
