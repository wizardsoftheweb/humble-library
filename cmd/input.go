package wotwhb

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	cowsay "github.com/Code-Hex/Neo-cowsay"
)

func buildPrompt(prompt string) string {
	say, err := cowsay.Say(
		cowsay.Phrase(prompt),
		cowsay.Type("default"),
	)
	fatalCheck(err)
	return say
}

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s\n", buildPrompt(prompt))
	input, err := reader.ReadString('\n')
	fatalCheck(err)
	return strings.TrimSpace(input)
}
