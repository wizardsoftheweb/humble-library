package wotwhb

import (
	. "gopkg.in/check.v1"
)

type LoggerSuite struct{}

var _ = Suite(&LoggerSuite{})

func (s *LoggerSuite) TestLoggerFormatter(c *C) {
	c.Assert(Logger.Formatter, Equals, formatter)
}
