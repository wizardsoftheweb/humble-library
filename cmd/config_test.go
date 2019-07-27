package wotwhb

import (
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type ConfigSuite struct {
	BaseSuite
	cookieFile string
}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) SetUpTest(c *C) {
	s.cookieFile = filepath.Join(ConfigDirectoryFlagValue, cookieFileBasename)
}

func (s *ConfigSuite) TestBuildSession(c *C) {
	client, jar := buildSession()
	_, err := os.Stat(s.cookieFile)
	c.Assert(err, NotNil)
	c.Assert(len(jar.AllCookies()), Equals, 0)
	c.Assert(client.Jar, Equals, jar)
}
