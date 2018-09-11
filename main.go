/* Cl is a tool for filtering data by columns on the command line.

Usage:

TODO

*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Flag definitions
var separator string
var useTabSeparator bool

// Use local variables and Reader/Writer interfaces so we can substitute these for testing
var stdin io.Reader
var stdout io.Writer
var stderr io.Writer
var exitFunc func(int)

func init() {
	flag.StringVar(&separator, "s", "", "a character or regex to split lines (default: whitespace)")
	flag.BoolVar(&useTabSeparator, "t", false, "use tabs as separator (alias of -s \\t)")

	stdin = os.Stdin
	stdout = os.Stdout
	stderr = os.Stderr
	exitFunc = os.Exit
}

func main() {
	exitFunc(run())
}

// run executes the command and returns a shell return code
func run() int {
	flag.Parse()

	// Figure out what columns (1-indexed) the user wants and validate them
	//

	columns := map[int64]struct{}{}
	for _, arg := range flag.Args() {
		c, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			fmt.Fprintf(stderr, "ERROR: failed to parse argument %q: %v\n", arg, err)
			return 1
		}
		if c < 1 {
			fmt.Fprintf(stderr, "ERROR: argument %q is invalid, column indexes must be positive numbers\n", c)
			return 1
		}
		columns[c] = struct{}{}
	}

	// Figure out how to separate columns
	//

	if separator != "" && useTabSeparator {
		fmt.Fprintf(stderr, "ERROR: you cannot use both -s and -t")
		return 1
	}

	if useTabSeparator {
		separator = "\t"
	}

	var separatorRegex *regexp.Regexp
	if separator != "" {
		var err error
		separatorRegex, err = regexp.Compile(separator)
		if err != nil {
			fmt.Fprintf(stderr, "ERROR: could not parse separator %q as a regular expression: %v\n", separator, err)
			return 1
		}
	}

	// Scan and split
	//

	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		var fields []string
		if separatorRegex == nil {
			fields = strings.Fields(scanner.Text())
		} else {
			fields = separatorRegex.Split(scanner.Text(), -1)
		}

		printedFirstColumn := false
		for i, f := range fields {
			if _, exists := columns[int64(i+1)]; exists {
				if printedFirstColumn {
					fmt.Fprintf(stdout, "\t")
				}
				printedFirstColumn = true
				fmt.Fprintf(stdout, f)
			}
		}
		fmt.Fprintf(stdout, "\n")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(stderr, "ERROR: failed to read input: %v\n", err)
		return 1
	}
	return 0
}
