package wotwhb

import (
	"github.com/spf13/cobra"
	. "gopkg.in/check.v1"
)

type WotwHbSuite struct {
	BaseSuite
	Command *cobra.Command
	Args    []string
}

const defaultVerbosityLevelFlag = 0

var _ = Suite(&WotwHbSuite{})

func (s *WotwHbSuite) SetUpTest(c *C) {
	VerbosityFlagValue = defaultVerbosityLevelFlag
	s.Command = &cobra.Command{}
	s.Args = []string{}
}

func (s *WotwHbSuite) TestInit(c *C) {
	c.Assert(VerbosityFlagValue, Equals, defaultVerbosityLevelFlag)
}

func (s *WotwHbSuite) TestPackageCmdPersistentPreRun(c *C) {
	c.Assert(Logger.Level, Equals, defaultLoggerLevel)
	for _, testLevel := range verbosityTestLevels {
		VerbosityFlagValue = testLevel.input
		PackageCmdPersistentPreRun(s.Command, s.Args)
		c.Assert(Logger.Level, Equals, testLevel.level)
	}
}

func (s *WotwHbSuite) TestPackageCmdRun(c *C) {
	c.Assert(PackageCmd.Version, Equals, PackageVersion)
	PackageCmdRun(s.Command, s.Args)
	c.Assert(PackageCmd.Version, Equals, PackageVersion)
}

func (s *WotwHbSuite) TestWotwHbExecute(c *C) {
	var oldPackageCmd = &cobra.Command{}
	*oldPackageCmd = *PackageCmd
	dummy := func(cmd *cobra.Command, args []string) {}
	PackageCmd.SilenceErrors = true
	PackageCmd.DisableFlagParsing = true
	PackageCmd.PersistentPreRun = dummy
	PackageCmd.PreRun = dummy
	PackageCmd.Run = dummy
	PackageCmd.PostRun = dummy
	PackageCmd.PersistentPostRun = dummy
	err := Execute()
	c.Assert(err, IsNil)
	*PackageCmd = *oldPackageCmd
}
