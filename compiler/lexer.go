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
	s    *bufio.Scanner
	next token
	sb   strings.Builder
}

func newLexer(s *bufio.Scanner) *lexer {
	s.Split(bufio.ScanRunes)
	s.Scan()
	return &lexer{s: s}
}

func (l *lexer) Peek() token {
	if l.next.T == Empty {
		l.advance()
	}
	return l.next
}

func (l *lexer) Pop() token {
	current := l.next
	l.advance()
	return current
}

// advance reads the next token in to the lexer
func (l *lexer) advance() {
	l.next = token{}
	for l.s.Text() == " " || l.s.Text() == "\t" {
		l.s.Scan()
	}
	scanAhead := true
	defer func() {
		if scanAhead {
			l.s.Scan()
		}
	}()
	switch {
	case l.s.Text() == "":
		l.next = token{}
	case l.s.Text() == "\n":
		l.next = token{T: NewLine}
	case strings.Contains("*+-/^=<>!", l.s.Text()):
		if l.s.Text() == "<" || l.s.Text() == ">" || l.s.Text() == "=" || l.s.Text() == "!" {
			s := l.s.Text()
			scanAhead = false
			l.s.Scan()
			if l.s.Text() == "=" {
				l.s.Scan()
				l.next = token{T: Op, V: s + "="}
			} else {
				l.next = token{T: Op, V: s}
			}
		} else {
			l.next = token{T: Op, V: l.s.Text()}
		}
	case l.s.Text() == "(":
		l.next = token{T: OpenBracket}
	case l.s.Text() == ")":
		l.next = token{T: CloseBracket}
	case strings.Contains("0123456789", l.s.Text()):
		scanAhead = false
		l.next = token{T: Number, V: l.parseNum()}
	default:
		r, _ := utf8.DecodeLastRuneInString(l.s.Text())
		if !(unicode.IsLetter(r) || r == '_') {
			l.next = token{T: Unknown, V: l.s.Text()}
			break
		}
		scanAhead = false
		l.next = token{T: Identifier, V: l.parseIdentifier()}
	}
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
