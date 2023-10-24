package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/JxBP/duden-check/internal/api"
)

func main() {
	if err := fallible_main(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
}

func fallible_main() error {
	var textBytes []byte
	var err error

	if len(os.Args) != 2 || os.Args[1] == "-" {
		textBytes, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading from stdin: %w", err)
		}
	} else {
		textBytes, err = os.ReadFile(os.Args[1])
		if err != nil {
			return fmt.Errorf("reading input file: %w", err)
		}
	}

	text := string(textBytes)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if err := checkText(line); err != nil {
			return err
		}
		fmt.Print("\n\n")
	}

	return nil
}

func checkText(text string) error {
	advices, err := api.FetchErrors(text)
	if err != nil {
		return err
	}

	fmt.Println(text)
	sort.Slice(advices, func(i, j int) bool {
		return advices[i].Offset < advices[j].Offset
	})

	{
		pos := 0

		for _, advice := range advices {
			indentLen := max(advice.Offset-pos, 0)
			indent := strings.Repeat(" ", indentLen)

			fmt.Printf("%s%s",
				indent,
				strings.Repeat("^", advice.Length),
			)
			pos += indentLen + advice.Length
		}
		fmt.Println()
	}

	for nAdvice := len(advices); nAdvice > 0; nAdvice-- {
		print_indent := func() int {
			print_partial_indent := func(pos int, advice api.SpellAdvice) {
				fmt.Print(strings.Repeat(" ", advice.Offset-pos))
			}

			pos := 0
			for _, advice := range advices[:nAdvice-1] {
				print_partial_indent(pos, advice)
				fmt.Print("|")
				pos = advice.Offset + 1
			}
			print_partial_indent(pos, advices[nAdvice-1])
			return pos
		}

		advice := advices[nAdvice-1]

		var proposalLinePrefix string
		if len(advice.Proposals) == 0 {
			proposalLinePrefix = "Keine Vorschläge"
		} else if len(advice.Proposals) == 1 {
			proposalLinePrefix = "Vorschlag: "
		} else {
			proposalLinePrefix = "Vorschläge: "
		}

		lines := []string{
			strings.Trim(advice.ShortMessage, "\r\n "),
			proposalLinePrefix + strings.Join(advice.Proposals, ", "),
		}

		print_indent()
		fmt.Println("|")

		width := longest(lines)

		border := func(s string) string { return fmt.Sprintf("+%s+", strings.Repeat(s, width+2)) }

		print_indent()
		fmt.Println(border("‾"))
		for _, line := range lines {
			print_indent()
			fmt.Printf("| %-[1]*[2]s |\n", width, line)
		}
		print_indent()
		fmt.Println(border("_"))
	}

	return nil
}

// Return the length of the longest string in the slice.
// If the slice is empty 0 is returned.
func longest(slice []string) int {
	rv := 0
	for _, s := range slice {
		rv = max(rv, len(s))
	}
	return rv
}
