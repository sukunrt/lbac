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
	expr(s, false)
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
	expr(s, false)
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

func expr(s *scanner, onlyOne bool) (nextToken string, err error) {
	skipSpaces(s)
	if s.Peek() == "" || s.Peek() == "\n" {
		return "", io.EOF
	}
	r, _ := utf8.DecodeRuneInString(s.Peek())
	switch {
	case r == '(':
		s.Pop()
		expr(s, false)
		if s.Pop() != ")" {
			return "", errors.New("invalid paranthesis")
		}
	case r == '+' || r == '-':
		op := s.Pop()
		expr(s, false)
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
			expr(s, false)
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
	if onlyOne {
		return "", nil
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
				emitOp("addq", "%rdi", "%rax")
			} else {
				emitOp("subq", "%rdi", "%rax")
			}
			push("%rax")
		case "*", "/":
			s.Pop()
			expr(s, true)
			pop("%rdi")
			pop("%rax")
			if op == "*" {
				emitOp("imulq", "%rdi", "%rax")
			} else {
				emitOp("cqto")
				emitOp("idivq", "%rdi")
			}
			push("%rax")
		case "\n", "", ")":
			return "", nil
		default:
			return "", fmt.Errorf("invalid character: %v", s.Peek())
		}
	}
}

func block(s *scanner) (nextToken string) {
	var err error
	for s.Peek() != "" {
		nextToken, err = expr(s, false)
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
