package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"
)

var lines []string

func emitln(s string) {
	lines = append(lines, s)
}

type scanner struct {
	s    *bufio.Scanner
	next string
}

func (s *scanner) peek() string {
	if s.next == "" {
		s.advance()
	}
	return s.next
}

func (s *scanner) pop() string {
	t := s.peek()
	s.advance()
	return t
}

func (s *scanner) advance() {
	s.s.Scan()
	s.next = s.s.Text()
}

func getNumber(s *scanner) (n int, err error) {
	skipSpaces(s)
	var started bool
	var breakCond string
	for {
		c := s.peek()
		if c == "" {
			breakCond = "EOF"
			break
		}
		started = true
		if c[0] >= '0' && c[0] <= '9' {
			n = n*10 + int(c[0]-'0')
			s.pop()
			continue
		} else {
			breakCond = "INVALID CHAR"
			break
		}
	}
	if !started {
		if breakCond == "EOF" {
			return 0, io.EOF
		} else {
			return 0, errors.New("NAN")
		}
	} else {
		return n, nil
	}
}

func getOPAdd(s *scanner) (op string, err error) {
	skipSpaces(s)
	c := s.peek()
	if c == "" {
		return "", io.EOF
	}
	switch c {
	case "+", "-":
		s.pop()
		return c, nil
	default:
		return "", fmt.Errorf("invalid token %s", c)
	}
}

func getOPMul(s *scanner) (op string, err error) {
	skipSpaces(s)
	c := s.peek()
	if c == "" {
		return "", io.EOF
	}
	switch c {
	case "*", "/":
		s.pop()
		return c, nil
	default:
		return "", fmt.Errorf("invalid token %s", c)
	}
}

func skipSpaces(s *scanner) {
	for {
		c := s.peek()
		if c == "" {
			return
		} else if unicode.IsSpace(rune(c[0])) {
			s.pop()
			continue
		} else {
			break
		}
	}
}

// getTerm parses a term and puts the result of the expression on top of stack
// A term is an expression with * or /
func getTerm(s *scanner) error {
	n, err := getNumber(s)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	emitln(fmt.Sprintf("\tmov\tx0, #%d", n))
	emitln(fmt.Sprintf("\tstr\tx0, [sp, -16]!"))
	for {
		op, err := getOPMul(s)
		if err != nil {
			return nil
		}
		n, err = getNumber(s)
		if err != nil {
			return fmt.Errorf("expected number: %w", err)
		}
		emitln(fmt.Sprintf("\tmov\tx1, #%d", n))
		emitln(fmt.Sprintf("\tldr\tx0, [sp], 16"))
		if op == "*" {
			emitln(fmt.Sprintf("\tmul\tx0, x1, x0"))
		} else {
			emitln(fmt.Sprintf("\tsdiv\tx0, x0, x1"))
		}
		emitln(fmt.Sprintf("\tstr\tx0, [sp, -16]!"))
	}
}

func parse(s *scanner) error {
	err := getTerm(s)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	for {
		op, err := getOPAdd(s)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		err = getTerm(s)
		if err != nil {
			return fmt.Errorf("expected number: %w", err)
		}
		emitln(fmt.Sprintf("\tldr\tx1, [sp], 16"))
		emitln(fmt.Sprintf("\tldr\tx0, [sp], 16"))
		if op == "+" {
			emitln(fmt.Sprintf("\tadd\tx0, x1, x0"))
		} else {
			emitln(fmt.Sprintf("\tsub\tx0, x0, x1"))
		}
		emitln(fmt.Sprintf("\tstr\tx0, [sp, -16]!"))
	}
}

func main() {
	emitln("\t.section\t__TEXT,__text,regular,pure_instructions")
	emitln("\t.globl _eval")
	emitln("\t.p2align\t2")
	emitln("_eval:")
	emitln("\t.cfi_startproc")
	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanRunes)
	ss := &scanner{s: s}
	if err := parse(ss); err != nil {
		fmt.Println(err)
		return
	}
	emitln("\tldr\tx0, [sp], 16")
	emitln("\tret")
	emitln("\t.cfi_endproc")
	emitln(".subsections_via_symbols")
	for _, s := range lines {
		fmt.Println(s)
	}
}
