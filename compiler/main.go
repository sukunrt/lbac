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

func newScanner(s *bufio.Scanner) *scanner {
	s.Split(bufio.ScanRunes)
	return &scanner{s: s, next: ""}
}

func (s *scanner) Peek() string {
	if s.next == "" {
		s.advance()
	}
	return s.next
}

func (s *scanner) Pop() string {
	t := s.Peek()
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
	var isEOF bool
	for {
		c := s.Peek()
		if c == "" {
			isEOF = true
			break
		}
		if c[0] >= '0' && c[0] <= '9' {
			started = true
			n = n*10 + int(c[0]-'0')
			s.Pop()
			continue
		}
		break
	}
	if !started {
		if isEOF {
			return 0, io.EOF
		} else {
			return 0, errors.New("NAN")
		}
	}
	return n, nil
}

func getOPAdd(s *scanner) (op string, err error) {
	skipSpaces(s)
	c := s.Peek()
	if c == "" {
		return "", io.EOF
	}
	switch c {
	case "+", "-":
		s.Pop()
		return c, nil
	default:
		return "", fmt.Errorf("invalid token %s", c)
	}
}

func getOPMul(s *scanner) (op string, err error) {
	skipSpaces(s)
	c := s.Peek()
	if c == "" {
		return "", io.EOF
	}
	switch c {
	case "*", "/":
		s.Pop()
		return c, nil
	default:
		return "", fmt.Errorf("invalid token %s", c)
	}
}

func skipSpaces(s *scanner) {
	for {
		c := s.Peek()
		if c == "" {
			return
		} else if unicode.IsSpace(rune(c[0])) {
			s.Pop()
			continue
		} else {
			return
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
	emitln(fmt.Sprintf(`	pushq	$%d`, n))
	for {
		op, err := getOPMul(s)
		if err != nil {
			return nil
		}
		n, err = getNumber(s)
		if err != nil {
			return fmt.Errorf("expected number: %w", err)
		}
		emitln(fmt.Sprintf(`	movq	$%d, %%rdi`, n))
		emitln(`	popq	%rax`)
		emitln(`	cqto`)
		if op == "*" {
			emitln(`	imulq	%rdi, %rax`)
		} else {
			emitln(`	idivq	%rdi`)
		}
		emitln(`	pushq	%rax`)
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
		emitln(`	popq	%rdi`)
		emitln(`	popq	%rax`)
		if op == "+" {
			emitln(`	addq 	%rdi, %rax`)
		} else {
			emitln(`	subq	%rdi, %rax`)
		}
		emitln(`	pushq %rax`)
	}
}

func main() {
	emitln(`	.text
	.file	"sm.c"
	.globl	eval                            # -- Begin function eval
	.p2align	4, 0x90
	.type	eval,@function
eval:`)
	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanRunes)
	ss := &scanner{s: s}
	if err := parse(ss); err != nil {
		fmt.Println(err)
		return
	}
	emitln(`	popq %rax`)
	emitln(`	retq
.Lfunc_end0:
	.size	eval, .Lfunc_end0-eval
					# -- End function
	.ident	"clang version 16.0.6"
	.section	".note.GNU-stack","",@progbits
	.addrsig
	`)
	for _, s := range lines {
		fmt.Println(s)
	}
}
