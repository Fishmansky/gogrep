package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var caseMap = map[byte]byte{
	'a': 'A',
	'A': 'a',
	'b': 'B',
	'B': 'b',
	'c': 'C',
	'C': 'c',
	'd': 'D',
	'D': 'd',
	'e': 'E',
	'E': 'e',
	'f': 'F',
	'F': 'f',
	'g': 'G',
	'G': 'g',
	'h': 'H',
	'H': 'h',
	'i': 'I',
	'I': 'i',
	'j': 'J',
	'J': 'j',
	'k': 'K',
	'K': 'k',
	'l': 'L',
	'L': 'l',
	'm': 'M',
	'M': 'm',
	'n': 'N',
	'N': 'n',
	'o': 'O',
	'O': 'o',
	'p': 'P',
	'P': 'p',
	'r': 'R',
	'R': 'r',
	's': 'S',
	'S': 's',
	't': 'T',
	'T': 't',
	'u': 'U',
	'U': 'u',
	'v': 'V',
	'V': 'v',
	'w': 'W',
	'W': 'w',
	'q': 'Q',
	'Q': 'q',
	'x': 'X',
	'X': 'x',
	'y': 'Y',
	'Y': 'y',
	'z': 'Z',
	'Z': 'z',
}

var usage = `usage: gogrep [-i] [-e pattern] [file ...]`

func pattern(s string, p string) bool {
	c := 0
	if caseIns {
		for i := 0; i < len(s); i++ {
			if s[i] == p[c] || s[i] == caseMap[p[c]] {
				if c == len(p)-1 {
					return true
				}
				c++
				continue
			}
			c = 0
		}
	} else {
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
	}
	return false
}

func patternRegex(s string, p string) bool {
	if caseIns {
		reg := regexp.MustCompile("(?i)" + p)
		if reg.MatchString(s) {
			return true
		}
		return false
	}
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

var caseIns bool
var regexStr string
var patternFileStr string

func main() {
	flag.StringVar(&regexStr, "e", "", "find provided pattern")
	flag.StringVar(&patternFileStr, "f", "", "specify input file with patterns")
	flag.BoolVar(&caseIns, "i", false, "case-insensitive")
	flag.Parse()
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println(usage)
		return
	}
	if patternFileStr != "" {
		pf, err := os.Open(patternFileStr)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(pf)
		for scanner.Scan() {
			for _, file := range flag.Args() {
				f, err := os.Open(file)
				if err != nil {
					panic(err)
				}
				if regexStr != "" {
					checkWith(f, scanner.Text(), patternRegex)
				}
				checkWith(f, scanner.Text(), pattern)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "readinf standard input: %s", err)
		}
		return
	}
	for _, file := range flag.Args()[1:] {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		if regexStr != "" {
			checkWith(f, regexStr, patternRegex)
		}
		checkWith(f, flag.Args()[0], pattern)
	}
}
