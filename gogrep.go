package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var usage = `usage: gogrep [-i] [-e pattern] [file ...]`

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

func checkWith(f *os.File, str string, fn func(s string, p string) bool) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if fn(scanner.Text(), str) {
			fmt.Println(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "readinf standard input: %s", err)
	}
}

func main() {
	var caseIns bool
	var regexStr string
	flag.StringVar(&regexStr, "e", "", "find provided pattern")
	flag.BoolVar(&caseIns, "i", false, "case-insensitive")
	flag.Parse()
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(usage)
		return
	}
	for _, file := range flag.Args() {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		if regexStr != "" {
			checkWith(f, regexStr, patternRegex)
		}
		checkWith(f, args[0], pattern)
	}
}
