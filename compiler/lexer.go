package main

import (
	"bufio"
	"slices"
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
	Keyword
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
	n := l.s.Text()
	switch {
	case n == "":
		next = token{}
	case n == "\n":
		next = token{T: NewLine}
	case strings.Contains("=!<>", n):
		scanAhead = false
		l.s.Scan()
		if l.s.Text() == "=" {
			l.s.Scan()
			next = token{T: Op, V: n + "="}
		} else {
			next = token{T: Op, V: n}
		}
	case strings.Contains("*+-/^", n):
		next = token{T: Op, V: n}
	case n == "(":
		next = token{T: OpenBracket}
	case n == ")":
		next = token{T: CloseBracket}
	case strings.Contains("0123456789", n):
		scanAhead = false
		next = token{T: Number, V: l.parseNum()}
	default:
		r, _ := utf8.DecodeLastRuneInString(n)
		if !(unicode.IsLetter(r) || r == '_') {
			next = token{T: Unknown, V: n}
			break
		}
		scanAhead = false
		s := l.parseIdentifier()
		if isKeyword(s) {
			next = token{T: Keyword, V: s}
		} else {
			next = token{T: Identifier, V: s}
		}
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

func isKeyword(s string) bool {
	return slices.Contains([]string{"IF", "ELSE", "ENDIF", "WHILE", "ENDWHILE", "FN", "ENDFN", "CALL"}, s)
}
