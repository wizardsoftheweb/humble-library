package wotwhb

import (
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type MetaSuite struct {
	BaseSuite
	choices          StringFromChoices
	absoluteFileName string
}

var _ = Suite(&MetaSuite{})

func (s *MetaSuite) SetUpTest(c *C) {
	s.choices = CompletionFormat
	s.choices.value = completionFormatChoices[0]
	s.absoluteFileName = filepath.Join(s.WorkingDir, "completions")
}

func (s *MetaSuite) TestChoiceType(c *C) {
	c.Assert(s.choices.Type(), Equals, completionFormatType)
}

func (s *MetaSuite) TestChoiceString(c *C) {
	c.Assert(s.choices.String(), Equals, completionFormatChoices[0])
}

func (s *MetaSuite) TestChoiceSet(c *C) {
	err := s.choices.Set(completionFormatChoices[1])
	c.Assert(err, IsNil)
	c.Assert(s.choices.value, Equals, completionFormatChoices[1])
	err = s.choices.Set("qqq")
	c.Assert(err, NotNil)

}

func (s *MetaSuite) TestCompletionGenerationPath(c *C) {
	CompletionOutputFilePath = s.absoluteFileName
	CompletionActionRun(s.Command, s.Args)
	_, err := os.Stat(CompletionOutputFilePath)
	c.Assert(err, IsNil)
	CompletionOutputFilePath = "test"
	CompletionActionRun(s.Command, s.Args)
	_, err = os.Stat(CompletionOutputFilePath)
	c.Assert(err, IsNil)
	removePath, _ := filepath.Abs(
		filepath.Join(
			".",
			CompletionOutputFilePath,
		),
	)
	_ = os.Remove(removePath)
}

func (s *MetaSuite) TestCompletionGenerationChoices(c *C) {
	CompletionOutputFilePath = s.absoluteFileName
	for _, choice := range completionFormatChoices {
		_ = CompletionFormat.Set(choice)

		CompletionActionRun(s.Command, s.Args)
		_, err := os.Stat(CompletionOutputFilePath)
		c.Assert(err, IsNil)
		_ = os.Remove(CompletionOutputFilePath)
	}
}
