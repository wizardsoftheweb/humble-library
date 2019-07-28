package wotwhb

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var CompletionFormat = StringFromChoices{
	choices: completionFormatChoices,
}
var CompletionOutputFilePath string

func init() {
	CompletionAction.Flags().VarP(
		&CompletionFormat,
		"format",
		"f",
		fmt.Sprintf("must be one of these choices: %s", CompletionFormat.choices),
	)
	CompletionAction.Flags().StringVarP(
		&CompletionOutputFilePath,
		"output",
		"o",
		"_wotw-hb_completions",
		"The output filename",
	)
	MetaCmd.AddCommand(CompletionAction)

	PackageCmd.AddCommand(MetaCmd)
}

var MetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Meta features like autocompletion",
	Run:   HelpOnly,
}

type StringFromChoices struct {
	choices []string
	value   string
}

func (s *StringFromChoices) String() string {
	return s.value
}

func (s *StringFromChoices) Set(value string) error {
	for _, element := range s.choices {
		if value == element {
			s.value = element
			return nil
		}
	}
	return errors.New(fmt.Sprintf("%s was not in the list of acceptable choices: %s", value, s.choices))
}

func (s *StringFromChoices) Type() string {
	return completionFormatType
}

var CompletionAction = &cobra.Command{
	Use:   "completion",
	Short: "Generate autocompletion files",
	Run:   CompletionActionRun,
}

func CompletionActionRun(cmd *cobra.Command, args []string) {
	var fileName string
	if filepath.IsAbs(CompletionOutputFilePath) {
		fileName = CompletionOutputFilePath
	} else {
		fileName, _ = filepath.Abs(
			filepath.Join(
				".",
				CompletionOutputFilePath,
			),
		)
	}
	switch CompletionFormat.value {
	case "bash":
		err := cmd.GenBashCompletionFile(fileName)
		fatalCheck(err)
	case "ps":
		err := cmd.GenPowerShellCompletionFile(fileName)
		fatalCheck(err)
	default:
		err := cmd.GenZshCompletionFile(fileName)
		fatalCheck(err)
	}
}
