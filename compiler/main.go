package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var lines []string

func emitln(s string) {
	lines = append(lines, s)
}

func emitOp(op string, operands ...string) {
	sb := strings.Builder{}
	if strings.HasSuffix(op, ":") {
		sb.WriteString(op)
	} else {
		sb.WriteString("\t" + op)
	}
	for i, o := range operands {
		if i == 0 {
			sb.WriteString("\t")
		} else {
			sb.WriteString(", ")
		}
		sb.WriteString(o)
	}
	emitln(sb.String())
}

var sp int

var variables = map[string]int{}

func push(s string) {
	sp++
	emitOp("pushq", s)
}

func pop(s string) {
	sp--
	emitOp("popq", s)
}

var labelCnt int

func newLabel() string {
	defer func() { labelCnt++ }()
	return fmt.Sprintf("JL%d", labelCnt)
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
	skipSpaces(s)
	if s.Peek() == "\n" || s.Peek() == "" {
		return errors.New("empty expression in IF condition")
	}
	expr(s, -1)
	pop("%rax")
	emitOp("cmpq", "$0", "%rax")
	l := newLabel()
	emitOp("je", l)
	nt := block(s)
	switch nt {
	case "ELSE":
		el := newLabel()
		emitOp("jmp", el)
		emitOp(l + ":")
		nt := block(s)
		if nt != "ENDIF" {
			return fmt.Errorf("expected ENDIF got %s", nt)
		}
		emitOp(el + ":")
	case "ENDIF":
		emitOp(l + ":")
		return nil
	default:
		return fmt.Errorf("expected ENDIF or ELSE got %s", nt)
	}
	return nil
}

func whileStmt(s *scanner) error {
	skipSpaces(s)
	if s.Peek() == "\n" || s.Peek() == "" {
		return errors.New("empty expression in IF condition")
	}
	stL := newLabel()
	emitOp(stL + ":")
	expr(s, -1)
	pop("%rax")
	emitOp("cmpq", "$0", "%rax")
	l := newLabel()
	emitOp("je", l)
	nt := block(s)
	if nt != "ENDWHILE" {
		return fmt.Errorf("invalid end condition: expected ENDWHILE got: %s", nt)
	}
	emitOp("jmp", stL)
	emitOp(l + ":")
	return nil
}

func addOp() {
	pop("%rdi")
	pop("%rax")
	emitOp("addq", "%rdi", "%rax")
	push("%rax")
}

func subOp() {
	pop("%rdi")
	pop("%rax")
	emitOp("subq", "%rdi", "%rax")
	push("%rax")
}

func mulOp() {
	pop("%rdi")
	pop("%rax")
	emitOp("imulq", "%rdi", "%rax")
	push("%rax")
}

func divOp() {
	pop("%rdi")
	pop("%rax")
	emitOp("cqto")
	emitOp("idivq", "%rdi")
	push("%rax")
}

func expOp() {
	sl := newLabel()
	el := newLabel()
	pop("%rdi")
	pop("%rdx")
	emitOp("movq", "$1", "%rbx")
	emitOp("movq", "$1", "%rax")
	emitOp(sl + ":")
	emitOp("cmp", "$0", "%rdi")
	emitOp("je", el)
	emitOp("imulq", "%rdx", "%rax")
	emitOp("subq", "%rbx", "%rdi")
	emitOp("jmp", sl)
	emitOp(el + ":")
	push("%rax")
}

var bindingPower = map[string]int{
	"+": 10,
	"-": 10,
	"*": 20,
	"/": 20,
	"^": 30,
	// Terminating symbols
	")":  -100,
	"\n": -100,
	"":   -100,
}

func expr(s *scanner, power int) (nextToken string, err error) {
	skipSpaces(s)
	if s.Peek() == "" || s.Peek() == "\n" {
		return "", io.EOF
	}
	r, _ := utf8.DecodeRuneInString(s.Peek())
	switch {
	case r == '(':
		s.Pop()
		expr(s, -1)
		if s.Pop() != ")" {
			return "", errors.New("invalid paranthesis")
		}
	case r == '+' || r == '-':
		op := s.Pop()
		expr(s, power)
		if op == "-" {
			pop("%rdi")
			emitOp("movq", "$-1", "%rax")
			emitOp("imulq", "%rax", "%rdi")
			push("%rdi")
		}
	case unicode.IsNumber(r):
		n, err := getNumber(s)
		if err != nil {
			return "", err
		}
		push(fmt.Sprintf("$%d", n))
	case unicode.IsLetter(r):
		v, err := getVariable(s)
		if err != nil {
			return "", err
		}
		if v == "IF" {
			return "", ifStmt(s)
		}
		if v == "WHILE" {
			return "", whileStmt(s)
		}
		if v == "ENDIF" || v == "ELSE" || v == "ENDWHILE" {
			return v, nil
		}
		skipSpaces(s)
		if s.Peek() == "=" {
			s.Pop()
			expr(s, -1)
			if p, ok := variables[v]; ok {
				pop("%rax")
				emitOp("movq", "%rax", fmt.Sprintf("-%d(%%rbp)", 8*p))
			} else {
				variables[v] = sp
			}
			return "", nil
		} else {
			if _, ok := variables[v]; !ok {
				return "", fmt.Errorf("invalid variable %s", v)
			}
			push(fmt.Sprintf("-%d(%%rbp)", variables[v]*8))
		}
	default:
		return "", fmt.Errorf("invalid token %s", s.Peek())
	}
	for {
		skipSpaces(s)
		op := s.Peek()
		p, ok := bindingPower[op]
		if !ok {
			return "", fmt.Errorf("invalid operation: %s", op)
		}
		if p <= power {
			return "", nil
		}
		s.Pop()
		expr(s, p)
		switch op {
		case "+":
			addOp()
		case "-":
			subOp()
		case "*":
			mulOp()
		case "/":
			divOp()
		case "^":
			expOp()
		}
	}
}

func block(s *scanner) (nextToken string) {
	var err error
	for s.Peek() != "" {
		nextToken, err = expr(s, -1)
		if err != nil && err != io.EOF {
			fmt.Println("err:", err)
			return
		}
		if nextToken != "" {
			break
		}
		if s.Peek() == "\n" {
			s.Pop()
		}
	}
	return
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
	emitOp("movq", "%rbp", "%rsp")
	emitOp("popq", "%rbp")
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
