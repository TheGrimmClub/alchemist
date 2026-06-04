// Package prompt provides small helpers for asking the user questions
// on the command line.
package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// Line asks a question and returns a trimmed single-line answer.
func Line(question string) string {
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// LineDefault is like Line but returns def if the user enters nothing.
func LineDefault(question, def string) string {
	answer := Line(fmt.Sprintf("%s [%s]: ", question, def))
	if answer == "" {
		return def
	}
	return answer
}

// Confirm asks a yes/no question. The default choice is shown capitalized.
func Confirm(question string, defaultYes bool) bool {
	suffix := " [y/N]: "
	if defaultYes {
		suffix = " [Y/n]: "
	}
	answer := strings.ToLower(Line(question + suffix))
	if answer == "" {
		return defaultYes
	}
	return answer == "y" || answer == "yes"
}
