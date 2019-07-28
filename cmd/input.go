package wotwhb

import (
	"bufio"
	"io"
	"os"
	"strings"

	cowsay "github.com/Code-Hex/Neo-cowsay"
)

var inputReader io.Reader = os.Stdin

type CanPrint interface {
	Printf(format string, args ...interface{})
}

func buildPrompt(prompt string) string {
	say, err := cowsay.Say(
		cowsay.Phrase(prompt),
		cowsay.Type("default"),
	)
	fatalCheck(err)
	return say
}

func getInput(printer CanPrint, longPrompt, shortPrompt string) string {
	reader := bufio.NewReader(inputReader)
	printer.Printf("%s\n%s: ", buildPrompt(longPrompt), shortPrompt)
	input, err := reader.ReadString('\n')
	fatalCheck(err)
	return strings.TrimSpace(input)
}
