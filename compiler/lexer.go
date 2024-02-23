package main

import (
	"bufio"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Type int

const (
	Empty Type = iota
	NewLine
	Op
	OpenBracket
	CloseBracket
	Number
	Identifier
	Unknown
)

type token struct {
	T Type
	V string
}

type lexer struct {
	s     *bufio.Scanner
	sb    strings.Builder
	stack []token
}

func newLexer(s *bufio.Scanner) *lexer {
	s.Split(bufio.ScanRunes)
	s.Scan()
	return &lexer{s: s}
}

func (l *lexer) Peek() token {
	var v token
	if len(l.stack) == 0 {
		l.advance()
	}
	if len(l.stack) > 0 {
		v = l.stack[len(l.stack)-1]
	}
	return v
}

func (l *lexer) Pop() token {
	if len(l.stack) == 0 {
		l.advance()
	}
	var v token
	if len(l.stack) > 0 {
		v = l.stack[len(l.stack)-1]
		l.stack = l.stack[:len(l.stack)-1]
	}
	return v
}

func (l *lexer) Push(t token) {
	l.stack = append(l.stack, t)
}

// advance reads the next token in to the lexer
func (l *lexer) advance() {
	next := token{}
	for l.s.Text() == " " || l.s.Text() == "\t" {
		l.s.Scan()
	}
	scanAhead := true
	switch {
	case l.s.Text() == "":
		next = token{}
	case l.s.Text() == "\n":
		next = token{T: NewLine}
	case strings.Contains("*+-/^=<>!", l.s.Text()):
		if l.s.Text() == "<" || l.s.Text() == ">" || l.s.Text() == "=" || l.s.Text() == "!" {
			s := l.s.Text()
			scanAhead = false
			l.s.Scan()
			if l.s.Text() == "=" {
				l.s.Scan()
				next = token{T: Op, V: s + "="}
			} else {
				next = token{T: Op, V: s}
			}
		} else {
			next = token{T: Op, V: l.s.Text()}
		}
	case l.s.Text() == "(":
		next = token{T: OpenBracket}
	case l.s.Text() == ")":
		next = token{T: CloseBracket}
	case strings.Contains("0123456789", l.s.Text()):
		scanAhead = false
		next = token{T: Number, V: l.parseNum()}
	default:
		r, _ := utf8.DecodeLastRuneInString(l.s.Text())
		if !(unicode.IsLetter(r) || r == '_') {
			next = token{T: Unknown, V: l.s.Text()}
			break
		}
		scanAhead = false
		next = token{T: Identifier, V: l.parseIdentifier()}
	}
	if scanAhead {
		l.s.Scan()
	}
	l.stack = append(l.stack, next)
}

func (l *lexer) parseNum() string {
	l.sb.Reset()
	l.sb.WriteString(l.s.Text())
	for l.s.Scan() {
		c := l.s.Text()
		if c == "" || c[0] < '0' || c[0] > '9' {
			break
		}
		l.sb.WriteString(c)
	}
	return l.sb.String()

}

func (l *lexer) parseIdentifier() string {
	l.sb.Reset()
	l.sb.WriteString(l.s.Text())
	for l.s.Scan() {
		c := l.s.Text()
		if c == "" {
			break
		}
		r, _ := utf8.DecodeRuneInString(c)
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			break
		}
		l.sb.WriteString(c)
	}
	return l.sb.String()
}
