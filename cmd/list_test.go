package wotwhb

import (
	. "gopkg.in/check.v1"
)

type ListSuite struct {
	BaseSuite
	uniqueList *UniqueStringList
}

var _ = Suite(&ListSuite{})

func (s *ListSuite) SetUpTest(c *C) {
	s.uniqueList = NewUniqueStringList("one")
}

func (s *ListSuite) TestListCreation(c *C) {
	c.Assert(NewUniqueStringList().Size(), Equals, 0)
}

func (s *ListSuite) TestListContents(c *C) {
	c.Assert(s.uniqueList.Size(), Equals, 1)
}

func (s *ListSuite) TestAdd(c *C) {
	s.uniqueList.Add("one", "two", "three", "two", "one", "")
	c.Assert(s.uniqueList.Size(), Equals, 3)
}

func (s *ListSuite) TestContents(c *C) {
	c.Assert(s.uniqueList.Contents(), DeepEquals, []string{"one"})
}
