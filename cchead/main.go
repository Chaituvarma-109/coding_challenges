package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	numBytes := flag.Int64("c", 0, "Number of bytes to print")
	numLines := flag.Int64("n", 10, "Number of lines to print")
	flag.Parse()

	var inp []byte

	if fileNames := flag.Args(); len(fileNames) == 0 {
		for i := 0; i < 10; i++ {
			fmt.Scanln(&inp)
			fmt.Println(string(inp))
		}
	} else {
		if len(fileNames) > 1 {
			for _, file := range fileNames {
				fmt.Printf("==> %s <==\n", file)
				content := read_file(file)
				if *numBytes != 0 {
					readNBytes(content, *numBytes)
				} else {
					readNLines(content, *numLines)
				}
			}
		} else {
			content := read_file(fileNames[0])
			if *numBytes != 0 {
				readNBytes(content, *numBytes)
			} else {
				readNLines(content, *numLines)
			}
		}
	}
}

func read_file(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("err %s", err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("err %s", err)
	}

	return data
}

func readNBytes(content []byte, n int64) {
	r := strings.NewReader(string(content))
	lr := io.LimitReader(r, n)

	if _, err := io.Copy(os.Stdout, lr); err != nil {
		fmt.Printf("err %s", err)
	}

	fmt.Print("\n")
}

func readNLines(content []byte, n int64) {
	r := strings.NewReader(string(content))
	scanner := bufio.NewScanner(r)
	lineNumber := 1

	for scanner.Scan() {
		if lineNumber > int(n) {
			break
		}
		line := scanner.Text()
		fmt.Println(line)
		lineNumber++
	}
}
