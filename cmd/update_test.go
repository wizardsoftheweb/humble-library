package wotwhb

import (
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	. "gopkg.in/check.v1"
)

type UpdateSuite struct {
	BaseSuite
	Command *cobra.Command
	Args    []string
}

func (s *UpdateSuite) Printf(format string, args ...interface{}) {}

var _ = Suite(&UpdateSuite{})

func (s *UpdateSuite) SetUpTest(c *C) {
	s.Command = &cobra.Command{}
	s.Command.SetOut(ioutil.Discard)
	s.Command.SetErr(ioutil.Discard)
	s.Args = []string{}
	inputReader = strings.NewReader(inputTestResult + "\n")
	Logger.SetLevel(logrus.FatalLevel)
}

func (s *UpdateSuite) TearDownTest(c *C) {
	Logger.SetLevel(defaultLoggerLevel)
}

func (s *UpdateSuite) TestUpdateKeyListCmdRun(c *C) {
	c.Assert(
		func() {
			UpdateKeyListCmdRun(s.Command, s.Args)
		},
		PanicMatches,
		".*EOF.*",
	)
}

func (s *UpdateSuite) TestUpdateCmdRun(c *C) {
	c.Assert(
		func() {
			UpdateCmdRun(s.Command, s.Args)
		},
		Not(PanicMatches),
		".*",
	)
}
