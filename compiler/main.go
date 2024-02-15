package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

var lines []string

func emitln(s string) {
	lines = append(lines, s)
}

var sp int

var variables = map[string]int{}

func push(s string) {
	sp++
	emitln(fmt.Sprintf("pushq	%s", s))
}

func pop(s string) {
	sp--
	emitln(fmt.Sprintf("popq	%s", s))
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
		}
		return 0, errors.New("NAN")
	}
	return n, nil
}

func getVariable(s *scanner) (res string, err error) {
	var started bool
	var isEOF bool
	for {
		c := s.Peek()
		if c == "" {
			isEOF = true
			break
		}
		r, _ := utf8.DecodeRuneInString(c)
		if unicode.IsDigit(r) || unicode.IsLetter(r) || c == "_" {
			started = true
			res += c
			s.Pop()
			continue
		}
		break
	}
	if !started {
		if isEOF {
			return "", io.EOF
		}
		return "", errors.New("not a variable")
	}
	return res, nil
}

func getFactor(s *scanner) (err error) {
	skipSpaces(s)
	if s.Peek() == "" {
		return io.EOF
	}
	r, _ := utf8.DecodeRuneInString(s.Peek())
	switch {
	case r == '(':
		s.Pop()
		getExpression(s)
		if s.Pop() != ")" {
			return errors.New("invalid paranthesis")
		}
	case r == '+' || r == '-':
		op := s.Pop()
		getExpression(s)
		if op == "-" {
			emitln(`	popq	%rdi`)
			emitln(`	movq	$-1, %rax`)
			emitln(`	imulq	%rax, %rdi`)
			push("%rdi")
		}
		return nil
	case unicode.IsNumber(r):
		n, err := getNumber(s)
		if err != nil {
			return err
		}
		push(fmt.Sprintf("$%d", n))
		return nil
	case unicode.IsLetter(r):
		v, err := getVariable(s)
		if err != nil {
			return err
		}
		skipSpaces(s)
		if s.Peek() == "=" {
			s.Pop()
			getExpression(s)
			variables[v] = sp
		} else {
			if _, ok := variables[v]; !ok {
				return fmt.Errorf("invalid variable %s", v)
			}
			push(fmt.Sprintf("-%d(%%rbp)", variables[v]*8))
		}
		return nil
	}
	return nil
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
		} else if unicode.IsSpace(rune(c[0])) && c != "\n" {
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
	err := getFactor(s)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	for {
		op, err := getOPMul(s)
		if err != nil {
			return nil
		}
		err = getFactor(s)
		if err != nil {
			return fmt.Errorf("expected number: %w", err)
		}
		pop("%rdi")
		pop("%rax")
		emitln(`	cqto`)
		if op == "*" {
			emitln(`	imulq	%rdi, %rax`)
		} else {
			emitln(`	idivq	%rdi`)
		}
		push("%rax")
	}
}

func getExpression(s *scanner) error {
	skipSpaces(s)
	err := getTerm(s)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	for {
		if s.Peek() == "\n" {
			return nil
		}
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
		pop("%rdi")
		pop("%rax")
		if op == "+" {
			emitln(`	addq 	%rdi, %rax`)
		} else {
			emitln(`	subq	%rdi, %rax`)
		}
		push("%rax")
	}
}

func main() {
	emitln(`	.text
	.file	"sm.c"
	.globl	eval                            # -- Begin function eval
	.p2align	4, 0x90
	.type	eval,@function
eval:
	pushq	%rbp
	movq	%rsp, %rbp
	`)
	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanRunes)
	ss := &scanner{s: s}
	for ss.Peek() != "" {
		err := getExpression(ss)
		if err != nil {
			fmt.Println(err)
			return
		}
		if ss.Peek() == "\n" {
			ss.Pop()
		}
	}
	pop("%rax")
	emitln(`	movq %rbp, %rsp`)
	emitln(`	popq %rbp`)
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
