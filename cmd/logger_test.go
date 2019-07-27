package wotwhb

import (
	"github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

type LoggerSuite struct {
	BaseSuite
}

var _ = Suite(&LoggerSuite{})

const defaultLoggerLevel = logrus.DebugLevel

func (*LoggerSuite) SetUpTest(c *C) {
	Logger.SetLevel(defaultLoggerLevel)
}

var verbosityTestLevels = []struct {
	input int
	level logrus.Level
}{
	{-10, logrus.PanicLevel},
	{-2, logrus.PanicLevel},
	{-1, logrus.FatalLevel},
	{0, logrus.ErrorLevel},
	{1, logrus.WarnLevel},
	{2, logrus.InfoLevel},
	{3, logrus.TraceLevel},
	{4, logrus.DebugLevel},
	{20, logrus.DebugLevel},
}

func (s *LoggerSuite) TestLoggerFormatter(c *C) {
	c.Assert(Logger.Formatter, Equals, formatter)
}

func (s *LoggerSuite) TestSettingLogLevel(c *C) {
	c.Assert(Logger.Level, Equals, defaultLoggerLevel)
	for _, testLevel := range verbosityTestLevels {
		setLoggerLevel(testLevel.input)
		c.Assert(Logger.Level, Equals, testLevel.level)
	}
}

func (s *LoggerSuite) TestBootstrapping(c *C) {
	c.Assert(Logger.Level, Equals, defaultLoggerLevel)
	for _, testLevel := range verbosityTestLevels {
		BootstrapLogger(testLevel.input)
		c.Assert(Logger.Level, Equals, testLevel.level)
	}
}
