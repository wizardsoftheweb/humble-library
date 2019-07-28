package wotwhb

import (
	"strings"

	. "gopkg.in/check.v1"
)

type InputSuite struct {
	BaseSuite
}

func (s *InputSuite) Printf(format string, args ...interface{}) {}

var _ = Suite(&InputSuite{})

const promptTestResults = ` _______
< hello >
 -------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||`

const inputTestResult = "this is some input"

func (s *InputSuite) SetUpTest(c *C) {
	inputReader = strings.NewReader(inputTestResult + "\n")
}

func (s *InputSuite) TestBuildPrompt(c *C) {
	c.Assert(buildPrompt("hello"), Equals, promptTestResults)
}

func (s *InputSuite) TestGetInput(c *C) {
	c.Assert(getInput(s, " ", " "), Equals, inputTestResult)
}
