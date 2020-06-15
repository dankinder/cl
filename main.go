/* Cl is a tool for filtering data by columns on the command line.

Usage: cl [options...] <column_indexes...>

Examples:

	Filter a simple table of data for the second column:
	$ echo "1 2 3
	> 4 5 6
	> 7 8 9" | cl 2
	2
	5
	8

	Grab a list of process IDs from ps, ignoring the header row (-i):
	$ ps | cl 1 -i
	7958
	29855

	Grab the third column of output when there may be spaces, in values, and tabs are the separator (-t):
	$ netstat | cl 3 -t

	Or the first 2 columns:
	$ netstat | cl 1 2 -t

	Or the first 4 columns (in bash):
	$ netstat | cl {1..4} -t

Options:
  -i    ignore the header row (first row)
  -s string
        a character or regex to split lines (default: whitespace)
  -t    use tabs as separator (alias of -s \t)
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
var ignoreHeaderRow bool

// Use local variables and Reader/Writer interfaces so we can substitute these for testing
var stdin io.Reader
var stdout io.Writer
var stderr io.Writer
var exitFunc func(int)

func init() {
	flag.StringVar(&separator, "s", "", "a character or regex to split lines (default: whitespace)")
	flag.BoolVar(&useTabSeparator, "t", false, "use tabs as separator (alias of -s \\t)")
	flag.BoolVar(&ignoreHeaderRow, "i", false, "ignore the header row (first row)")

	flag.Usage = func() {
		fmt.Printf(`Usage: cl [options...] <column_indexes...>

cl is a tool for filtering data by columns on the command line.

Examples:

	Filter a simple table of data for the second column:
	$ echo "1 2 3
	> 4 5 6
	> 7 8 9" | cl 2
	2
	5
	8

	Grab a list of process IDs from ps, ignoring the header row (-i):
	$ ps | cl 1 -i
	7958
	29855

	Grab the third column of output when there may be spaces, in values, and tabs are the separator (-t):
	$ netstat | cl 3 -t

	Or the first 2 columns:
	$ netstat | cl 1 2 -t

	Or the first 4 columns (in bash):
	$ netstat | cl {1..4} -t

Options:
`)
		flag.PrintDefaults()
	}

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

	if len(columns) == 0 {
		flag.Usage()
		os.Exit(1)
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
	firstRow := true
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		if firstRow && ignoreHeaderRow {
			firstRow = false
			continue
		}

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
					fmt.Fprint(stdout, "\t")
				}
				printedFirstColumn = true
				fmt.Fprint(stdout, f)
			}
		}
		fmt.Fprint(stdout, "\n")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(stderr, "ERROR: failed to read input: %v\n", err)
		return 1
	}
	return 0
}
