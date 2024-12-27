package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

func fileStats(file []byte) (int, int, int, int) {
	fb := len(file)
	charCount := len(bytes.Runes(file))
	wordCount := len(bytes.Fields(file))
	lines := bytes.Count(file, []byte("\n"))

	return fb, lines, wordCount, charCount
}

func main() {
	cVal := flag.Bool("c", false, "to print number of bytes in a file.")
	lVal := flag.Bool("l", false, "to print number of lines in a file.")
	wVal := flag.Bool("w", false, "to print number of words in a file.")
	mVal := flag.Bool("m", false, "to print number of characters in a file.")

	flag.Parse()

	var file []byte
	args := flag.Args()

	if len(args) > 0 {
		f, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		file = f
	} else {
		fin, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		file = fin
	}

	switch {
	case *cVal:
		bytes, _, _, _ := fileStats(file)
		fmt.Printf("%d %s\n", bytes, flag.Arg(0))
	case *lVal:
		_, lines, _, _ := fileStats(file)
		fmt.Printf("%d %s\n", lines, flag.Arg(0))
	case *wVal:
		_, _, words, _ := fileStats(file)
		fmt.Printf("%d %s\n", words, flag.Arg(0))
	case *mVal:
		_, _, _, chars := fileStats(file)
		fmt.Printf("%d %s\n", chars, flag.Arg(0))
	default:
		fb, lines, wordCount, _ := fileStats(file)
		fmt.Printf(" %d %d %d %s\n", lines, wordCount, fb, flag.Arg(0))
	}
}
