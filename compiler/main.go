package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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

func ifStmt(l *lexer) error {
	if l.Peek().T == NewLine || l.Peek().T == Empty {
		return errors.New("empty expression in IF condition")
	}
	expr(l, -1)
	pop("%rax")
	emitOp("cmpq", "$0", "%rax")
	elseL := newLabel()
	emitOp("je", elseL)
	nt := block(l)
	switch nt {
	case "ELSE":
		endL := newLabel()
		emitOp("jmp", endL)
		emitOp(elseL + ":")
		nt := block(l)
		if nt != "ENDIF" {
			return fmt.Errorf("expected ENDIF got %s", nt)
		}
		emitOp(endL + ":")
	case "ENDIF":
		emitOp(elseL + ":")
		return nil
	default:
		return fmt.Errorf("expected ENDIF or ELSE got %s", nt)
	}
	return nil
}

func whileStmt(l *lexer) error {
	if l.Peek().T == NewLine || l.Peek().T == Empty {
		return errors.New("empty expression in IF condition")
	}
	stL := newLabel()
	emitOp(stL + ":")
	expr(l, -1)
	pop("%rax")
	emitOp("cmpq", "$0", "%rax")
	el := newLabel()
	emitOp("je", el)
	nt := block(l)
	if nt != "ENDWHILE" {
		return fmt.Errorf("invalid end condition: expected ENDWHILE got: %s", nt)
	}
	emitOp("jmp", stL)
	emitOp(el + ":")
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

func expr(l *lexer, power int) (nextToken string, err error) {
	if l.Peek().T == Empty || l.Peek().T == NewLine {
		return "", io.EOF
	}
	switch l.Peek().T {
	case OpenBracket:
		l.Pop()
		expr(l, -1)
		if l.Pop().T != CloseBracket {
			return "", errors.New("invalid paranthesis")
		}
	case Op:
		op := l.Pop()
		if op.V != "+" && op.V != "-" {
			return "", fmt.Errorf("invalid op: %+v", op)
		}
		expr(l, power)
		if op.V == "-" {
			pop("%rdi")
			emitOp("movq", "$-1", "%rax")
			emitOp("imulq", "%rax", "%rdi")
			push("%rdi")
		}
	case Number:
		push(fmt.Sprintf("$%s", l.Peek().V))
		l.Pop()
	case Identifier:
		v := l.Peek().V
		l.Pop()
		if v == "IF" {
			return "", ifStmt(l)
		}
		if v == "WHILE" {
			return "", whileStmt(l)
		}
		if v == "ENDIF" || v == "ELSE" || v == "ENDWHILE" {
			return v, nil
		}
		if l.Peek().T == Op && l.Peek().V == "=" {
			l.Pop()
			expr(l, -1)
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
		return "", fmt.Errorf("invalid token %v", l.Peek())
	}
	for {
		op := l.Peek()
		if op.T == Empty || op.T == NewLine || op.T == CloseBracket {
			return "", nil
		}
		if op.T != Op {
			return "", fmt.Errorf("invalid operation: %+v", op)
		}
		p, ok := bindingPower[op.V]
		if !ok {
			return "", fmt.Errorf("invalid operation: %+v", op)
		}
		if p <= power {
			return "", nil
		}
		l.Pop()
		expr(l, p)
		switch op.V {
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
		default:
			return "", fmt.Errorf("invalid operation: %+v", op)
		}
	}
}

func block(l *lexer) (nextToken string) {
	var err error
	for l.Peek().T != Empty {
		nextToken, err = expr(l, -1)
		if err != nil && err != io.EOF {
			fmt.Println("err:", err)
			return
		}
		if nextToken != "" {
			break
		}
		if l.Peek().T == NewLine {
			l.Pop()
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
	l := newLexer(s)
	block(l)
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
