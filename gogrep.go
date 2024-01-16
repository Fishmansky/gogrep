package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var usage = `usage: gogrep [-e pattern] [file ...]`

func pattern(s string, p string) bool {
	c := 0
	for i := 0; i < len(s); i++ {
		if s[i] == p[c] {
			if c == len(p)-1 {
				return true
			}
			c++
			continue
		}
		c = 0
	}
	return false
}

func patternRegex(s string, p string) bool {
	reg := regexp.MustCompile(p)
	if reg.MatchString(s) {
		return true
	}
	return false
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
			if patternRegex(scanner.Text(), regexStr) {
				fmt.Println(scanner.Text())
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "readinf standard input: %s", err)
		}
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if pattern(scanner.Text(), args[0]) {
			fmt.Println(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "readinf standard input: %s", err)
	}

}
