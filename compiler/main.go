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

func ifStmt(s *scanner) error {
	return nil
}

func expr(s *scanner, onlyOne bool) error {
	skipSpaces(s)
	if s.Peek() == "" || s.Peek() == "\n" {
		return io.EOF
	}
	r, _ := utf8.DecodeRuneInString(s.Peek())
	switch {
	case r == '(':
		s.Pop()
		expr(s, false)
		if s.Pop() != ")" {
			return errors.New("invalid paranthesis")
		}
	case r == '+' || r == '-':
		op := s.Pop()
		expr(s, false)
		if op == "-" {
			pop("%rdi")
			emitln(`	movq	$-1, %rax`)
			emitln(`	imulq	%rax, %rdi`)
			push("%rdi")
		}
	case unicode.IsNumber(r):
		n, err := getNumber(s)
		if err != nil {
			return err
		}
		push(fmt.Sprintf("$%d", n))
	case unicode.IsLetter(r):
		v, err := getVariable(s)
		if err != nil {
			return err
		}
		if v == "IF" {
			return ifStmt(s)
		}
		skipSpaces(s)
		if s.Peek() == "=" {
			s.Pop()
			expr(s, false)
			variables[v] = sp
			return nil
		} else {
			if _, ok := variables[v]; !ok {
				return fmt.Errorf("invalid variable %s", v)
			}
			push(fmt.Sprintf("-%d(%%rbp)", variables[v]*8))
		}
	}
	if onlyOne {
		return nil
	}
	for {
		skipSpaces(s)
		op := s.Peek()
		switch op {
		case "+", "-":
			s.Pop()
			expr(s, false)
			pop("%rdi")
			pop("%rax")
			if op == "+" {
				emitln(`	addq 	%rdi, %rax`)
			} else {
				emitln(`	subq	%rdi, %rax`)
			}
			push("%rax")
		case "*", "/":
			s.Pop()
			expr(s, true)
			pop("%rdi")
			pop("%rax")
			if op == "*" {
				emitln(`	imulq 	%rdi, %rax`)
			} else {
				emitln(`cqto`)
				emitln(`	idivq	%rdi`)
			}
			push("%rax")
		case "\n", "", ")":
			return nil
		default:
			return fmt.Errorf("invalid character: %v", s.Peek())
		}
	}
}

func block(s *scanner) {
	st := sp
	for s.Peek() != "" {
		err := expr(s, false)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		if s.Peek() == "\n" {
			s.Pop()
		}
	}
	pop("%rax")
	emitln(`	movq	%rbp, %rdi`)
	emitln(fmt.Sprintf(`	movq	$%d, %%rbx`, 8*st))
	emitln(`	subq	%rbx, %rdi`)
	emitln(`	movq	%rdi, %rsp`)
	push("%rax")
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
	block(ss)
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
