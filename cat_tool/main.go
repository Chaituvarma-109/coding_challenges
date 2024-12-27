package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var printLineNums, printNonBlankLineNums bool

func main() {
	flag.BoolVar(&printLineNums, "n", false, "to print line nums")
	flag.BoolVar(&printNonBlankLineNums, "b", false, "to print non-blank line nums")
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		args = append(args, "-")
	}

	for _, arg := range args {
		var file *os.File
		var err error

		if arg == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(arg)
			if err != nil {
				fmt.Printf("err %s", err)
			}
			defer file.Close()
		}

		lineNum := 0
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()

			if printLineNums {
				lineNum++
				fmt.Printf("%d. %s\n", lineNum, line)
			} else if printNonBlankLineNums {
				if len(strings.TrimRight(line, "\r\n")) > 0 {
					lineNum++
					fmt.Printf("%d. %s\n", lineNum, line)
				}
			} else {
				fmt.Println(line)
			}
		}
	}
}
