package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var usage = `usage: gogrep [-e pattern] [file ...]`

func match(s string, p string) {
	reg := regexp.MustCompile(p)
	fmt.Println(reg.FindAllString(s, -1))
}

func main() {
	var regexStr string
	flag.StringVar(&regexStr, "e", "", "find provided pattern")
	flag.Parse()
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(usage)
		return
	}
	f, err := os.Open(flag.Args()[len(flag.Args())-1])
	if err != nil {
		panic(err)
	}
	if regexStr != "" {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			match(scanner.Text(), regexStr)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "readinf standard input: %s", err)
		}
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "readinf standard input: %s", err)
	}

}
