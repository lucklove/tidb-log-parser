package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/scanner"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Println(strings.Join(parseLine(sc.Text()), ""))
	}
}

func parseLine(line string) []string {
	s := scanner.Scanner{}
	s.Init(strings.NewReader(line))

	s.Error = func(s *scanner.Scanner, msg string) {}
	s.IsIdentRune = func(ch rune, i int) bool {
		r := !strings.ContainsRune(`[=]"`, ch) && ch != scanner.EOF
		return r
	}

	xs := []string{}
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		xs = append(xs, s.TokenText())
	}
	return xs
}
