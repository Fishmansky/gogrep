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

var usage = `usage: gogrep [-i] [-h] [-H] [-e pattern | -f pattern file] [file ...]`

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
	} else {
		reg := regexp.MustCompile(p)
		if reg.MatchString(s) {
			return true
		}
	}
	return false
}

func checkWith(f *os.File, str string, fn func(s string, p string) bool) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if fn(scanner.Text(), str) {
			if hideFilename {
				fmt.Printf("%s\n", scanner.Text())
			} else {
				fmt.Printf("%s:%s\n", f.Name(), scanner.Text())
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
	}
}

var caseIns bool
var regexStr string
var patternFileStr string
var hideFilename bool
var showFilename bool

func main() {
	flag.StringVar(&regexStr, "e", "", "find provided pattern")
	flag.BoolVar(&hideFilename, "h", false, "omit filenames in output")
	flag.BoolVar(&showFilename, "H", false, "print the file name for reach match")
	flag.StringVar(&patternFileStr, "f", "", "specify input file with patterns")
	flag.BoolVar(&caseIns, "i", false, "case-insensitive")
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println(usage)
		return
	}
	files := []string{}
	ptrn := ""
	if patternFileStr == "" && regexStr == "" {
		files = flag.Args()[1:]
		ptrn = flag.Args()[0]
	} else {
		files = flag.Args()
	}
	if len(files) == 1 {
		f, err := os.Open(files[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
		}
		if regexStr != "" {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if patternRegex(scanner.Text(), regexStr) {
					if showFilename {
						fmt.Printf("%s:%s\n", f.Name(), scanner.Text())
					} else {
						fmt.Printf("%s\n", scanner.Text())
					}
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
			}
		}
		if patternFileStr != "" {
			pf, err := os.Open(patternFileStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
			}
			scanner := bufio.NewScanner(pf)
			for scanner.Scan() {
				for _, file := range files {
					f, err := os.Open(file)
					if err != nil {
						fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
					}
					scanner := bufio.NewScanner(f)
					for scanner.Scan() {
						if pattern(scanner.Text(), ptrn) {
							if showFilename {
								fmt.Printf("%s:%s\n", f.Name(), scanner.Text())
							} else {
								fmt.Printf("%s\n", scanner.Text())
							}
						}
					}
					if err := scanner.Err(); err != nil {
						fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
					}
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
			}
			return
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if pattern(scanner.Text(), ptrn) {
				if showFilename {
					fmt.Printf("%s:%s\n", f.Name(), scanner.Text())
				} else {
					fmt.Printf("%s\n", scanner.Text())
				}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
		}
	} else {
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
			}
			if regexStr != "" {
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					if patternRegex(scanner.Text(), regexStr) {
						if hideFilename {
							fmt.Printf("%s\n", scanner.Text())
						} else {
							fmt.Printf("%s:%s\n", f.Name(), scanner.Text())
						}
					}
				}
				if err := scanner.Err(); err != nil {
					fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
				}
			}
			if patternFileStr != "" {
				pf, err := os.Open(patternFileStr)
				if err != nil {
					fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
				}
				scanner := bufio.NewScanner(pf)
				for scanner.Scan() {
					for _, file := range files {
						f, err := os.Open(file)
						if err != nil {
							fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
						}
						scanner := bufio.NewScanner(f)
						for scanner.Scan() {
							if pattern(scanner.Text(), ptrn) {
								if hideFilename {
									fmt.Printf("%s\n", scanner.Text())
								} else {
									fmt.Printf("%s:%s\n", f.Name(), scanner.Text())
								}
							}
						}
						if err := scanner.Err(); err != nil {
							fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
						}
					}
				}
				if err := scanner.Err(); err != nil {
					fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
				}
				return
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if pattern(scanner.Text(), ptrn) {
					if hideFilename {
						fmt.Printf("%s\n", scanner.Text())
					} else {
						fmt.Printf("%s:%s\n", f.Name(), scanner.Text())
					}
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "gogrep: %s\n", err)
			}
		}
	}

}
