// Copyright (c) 2013 Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// CAUTION: If this file is a Go source file (*.go), it was generated
// automatically by '$ golex' from a *.l file - DO NOT EDIT in that case!

// Package scanner implements a scanner for Go source text. It takes a
// []byte as source which can then be tokenized through repeated calls to the
// Scan method.
package scanner

import (
	"errors"
	"fmt"
	"go/token"
	"unicode"
	"unicode/utf8"
)

// A Scanner holds the scanner's internal state while processing a given text.
type Scanner struct {
	Col      int     // Starting column of the last scanned token.
	Errors   []error // List of accumulated errors.
	Fname    string  // File name (reported) of the scanned source.
	Line     int     // Starting line of the last scanned token.
	NCol     int     // Starting column (reported) for the next scanned token.
	NLine    int     // Starting line (reported) for the next scanned token.
	c        int
	i        int
	i0       int
	sc       int
	semi     bool
	src      []byte
	val      []byte
	lcomment string
}

// New returns a newly created Scanner.
func New(src []byte) (s *Scanner) {
	if len(src) > 2 && src[0] == 0xEF && src[1] == 0xBB && src[2] == 0xBF {
		src = src[3:]
	}
	s = &Scanner{
		src:   src,
		NLine: 1,
		NCol:  0,
	}
	s.next()
	return
}

func (s *Scanner) next() int {
	if s.c != 0 {
		s.val = append(s.val, byte(s.c))
	}
	s.c = 0
	if s.i < len(s.src) {
		s.c = int(s.src[s.i])
		s.i++
	}
	switch s.c {
	case '\n':
		s.NLine++
		s.NCol = 0
	default:
		s.NCol++
	}
	return s.c
}

// Pos returns the starting offset of the last scanned token.
func (s *Scanner) Pos() int {
	return s.i0
}

func (s *Scanner) err(format string, arg ...interface{}) {
	err := fmt.Errorf(fmt.Sprintf("%s:%d:%d ", s.Fname, s.Line, s.Col)+format, arg...)
	s.Errors = append(s.Errors, err)
}

// Error implements yyLexer.
func (s *Scanner) Error(msg string) {
	if msg == "syntax error" { // yacc
		s.err(msg)
		return
	}

	s.Errors = append(s.Errors, errors.New(msg))
}

// Scan scans the next token and returns the token and its value if applicable.
// The source end is indicated by token.EOF.
//
// If the returned token is a literal (token.IDENT, token.INT, token.FLOAT,
// token.IMAG, token.CHAR, token.STRING) or token.COMMENT, lval has has the
// corresponding value - not the string representation of the value. However,
// numeric literals (token.INT, token.FLOAT, token.IMAG) which overflow the
// corresponding Go predeclared types are returned as string.
//
// If the returned token is token.ILLEGAL, the literal string is the offending
// character or number/string/char literal.
func (s *Scanner) Scan() (tok token.Token, lval interface{}) {
	for {
		tok, lval = s.scan()
		if tok != token.ILLEGAL || lval.(string) != "\n" {
			return
		}
	}
}

// ScanSemis is like Scan but inserts semicolons as specified in [0].
//
//	[0]: http://golang.org/ref/spec#Semicolons
func (s *Scanner) ScanSemis() (tok token.Token, lval interface{}) {
	if s.lcomment != "" {
		tok, lval, s.lcomment = token.COMMENT, s.lcomment, ""
		return
	}

again:
	tok, lval = s.scan()
	switch tok {
	case token.EOF:
		if s.semi {
			tok, lval = token.SEMICOLON, "\n"
			s.semi = false
		}
	case token.ILLEGAL:
		if lval.(string) != "\n" {
			break
		}

		if !s.semi {
			goto again
		}

		tok, lval = token.SEMICOLON, "\n"
		s.semi = false
	case token.BREAK, token.CHAR, token.CONTINUE, token.DEC,
		token.FALLTHROUGH, token.FLOAT, token.IDENT, token.IMAG,
		token.INC, token.INT, token.RBRACE, token.RBRACK, token.RETURN,
		token.RPAREN, token.STRING:
		s.semi = true
	case token.COMMENT:
		if !s.semi {
			break
		}

		c := lval.(string)
		if c[:2] == "//" {
			s.lcomment = c
			tok, lval = token.SEMICOLON, "\n"
			s.semi = false
			break
		}

		src := s.src[s.i0-1+len(c):]
	search:
		for i := range src {
			switch src[i] {
			case ' ', '\t':
				//nop
			case '\n':
				s.lcomment = c
				tok, lval = token.SEMICOLON, "\n"
				s.semi = false
				break search
			default:
				break search
			}
		}
	default:
		s.semi = false
	}
	return
}

func (s *Scanner) scan() (tok token.Token, lval interface{}) {
	//defer func() { fmt.Printf("%s(%d) %v\n", tok, int(tok), lval) }()
	const (
		INITIAL = iota
		S1
		S2
	)

	c0, c := s.c, s.c

yystate0:

	s.val = s.val[:0]
	s.i0, s.Line, s.Col, c0 = s.i, s.NLine, s.NCol, c

	switch yyt := s.sc; yyt {
	default:
		panic(fmt.Errorf(`invalid start condition %d`, yyt))
	case 0: // start condition: INITIAL
		goto yystart1
	case 1: // start condition: S1
		goto yystart193
	case 2: // start condition: S2
		goto yystart197
	}

	goto yystate1 // silence unused label error
yystate1:
	c = s.next()
yystart1:
	switch {
	default:
		goto yyabort
	case c == '!':
		goto yystate4
	case c == '"':
		goto yystate6
	case c == '%':
		goto yystate7
	case c == '&':
		goto yystate9
	case c == '(':
		goto yystate20
	case c == ')':
		goto yystate21
	case c == '*':
		goto yystate22
	case c == '+':
		goto yystate24
	case c == ',':
		goto yystate27
	case c == '-':
		goto yystate28
	case c == '.':
		goto yystate31
	case c == '/':
		goto yystate39
	case c == '0':
		goto yystate45
	case c == ':':
		goto yystate52
	case c == ';':
		goto yystate54
	case c == '<':
		goto yystate55
	case c == '=':
		goto yystate60
	case c == '>':
		goto yystate62
	case c == '[':
		goto yystate67
	case c == '\'':
		goto yystate14
	case c == '\t' || c == '\r' || c == ' ':
		goto yystate3
	case c == '\x00':
		goto yystate2
	case c == ']':
		goto yystate68
	case c == '^':
		goto yystate69
	case c == '`':
		goto yystate71
	case c == 'b':
		goto yystate72
	case c == 'c':
		goto yystate77
	case c == 'd':
		goto yystate93
	case c == 'e':
		goto yystate102
	case c == 'f':
		goto yystate106
	case c == 'g':
		goto yystate122
	case c == 'i':
		goto yystate126
	case c == 'm':
		goto yystate141
	case c == 'p':
		goto yystate144
	case c == 'r':
		goto yystate151
	case c == 's':
		goto yystate161
	case c == 't':
		goto yystate177
	case c == 'v':
		goto yystate181
	case c == '{':
		goto yystate184
	case c == '|':
		goto yystate185
	case c == '}':
		goto yystate188
	case c >= '1' && c <= '9':
		goto yystate51
	case c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'h' || c >= 'j' && c <= 'l' || c == 'n' || c == 'o' || c == 'q' || c == 'u' || c >= 'w' && c <= 'z':
		goto yystate66
	case c >= 'Â' && c <= 'ß':
		goto yystate189
	case c >= 'à' && c <= 'ï':
		goto yystate191
	case c >= 'ð' && c <= 'ô':
		goto yystate192
	}

yystate2:
	c = s.next()
	goto yyrule1

yystate3:
	c = s.next()
	switch {
	default:
		goto yyrule2
	case c == '\t' || c == '\r' || c == ' ':
		goto yystate3
	}

yystate4:
	c = s.next()
	switch {
	default:
		goto yyrule5
	case c == '=':
		goto yystate5
	}

yystate5:
	c = s.next()
	goto yyrule6

yystate6:
	c = s.next()
	goto yyrule81

yystate7:
	c = s.next()
	switch {
	default:
		goto yyrule7
	case c == '=':
		goto yystate8
	}

yystate8:
	c = s.next()
	goto yyrule8

yystate9:
	c = s.next()
	switch {
	default:
		goto yyrule9
	case c == '&':
		goto yystate10
	case c == '=':
		goto yystate11
	case c == '^':
		goto yystate12
	}

yystate10:
	c = s.next()
	goto yyrule10

yystate11:
	c = s.next()
	goto yyrule11

yystate12:
	c = s.next()
	switch {
	default:
		goto yyrule12
	case c == '=':
		goto yystate13
	}

yystate13:
	c = s.next()
	goto yyrule13

yystate14:
	c = s.next()
	switch {
	default:
		goto yyrule84
	case c == '\'':
		goto yystate17
	case c == '\\':
		goto yystate18
	case c >= '\x01' && c <= '&' || c >= '(' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate15
	}

yystate15:
	c = s.next()
	switch {
	default:
		goto yyrule84
	case c == '\'':
		goto yystate16
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate15
	}

yystate16:
	c = s.next()
	goto yyrule85

yystate17:
	c = s.next()
	goto yyrule83

yystate18:
	c = s.next()
	switch {
	default:
		goto yyrule84
	case c == '\'':
		goto yystate19
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate15
	}

yystate19:
	c = s.next()
	switch {
	default:
		goto yyrule84
	case c == '\'':
		goto yystate16
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate15
	}

yystate20:
	c = s.next()
	goto yyrule14

yystate21:
	c = s.next()
	goto yyrule15

yystate22:
	c = s.next()
	switch {
	default:
		goto yyrule16
	case c == '=':
		goto yystate23
	}

yystate23:
	c = s.next()
	goto yyrule17

yystate24:
	c = s.next()
	switch {
	default:
		goto yyrule18
	case c == '+':
		goto yystate25
	case c == '=':
		goto yystate26
	}

yystate25:
	c = s.next()
	goto yyrule19

yystate26:
	c = s.next()
	goto yyrule20

yystate27:
	c = s.next()
	goto yyrule21

yystate28:
	c = s.next()
	switch {
	default:
		goto yyrule22
	case c == '-':
		goto yystate29
	case c == '=':
		goto yystate30
	}

yystate29:
	c = s.next()
	goto yyrule23

yystate30:
	c = s.next()
	goto yyrule24

yystate31:
	c = s.next()
	switch {
	default:
		goto yyrule25
	case c == '.':
		goto yystate32
	case c >= '0' && c <= '9':
		goto yystate34
	}

yystate32:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate33
	}

yystate33:
	c = s.next()
	goto yyrule26

yystate34:
	c = s.next()
	switch {
	default:
		goto yyrule80
	case c == 'E' || c == 'e':
		goto yystate35
	case c == 'i':
		goto yystate38
	case c >= '0' && c <= '9':
		goto yystate34
	}

yystate35:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '+' || c == '-':
		goto yystate36
	case c >= '0' && c <= '9':
		goto yystate37
	}

yystate36:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate37
	}

yystate37:
	c = s.next()
	switch {
	default:
		goto yyrule80
	case c == 'i':
		goto yystate38
	case c >= '0' && c <= '9':
		goto yystate37
	}

yystate38:
	c = s.next()
	goto yyrule79

yystate39:
	c = s.next()
	switch {
	default:
		goto yyrule27
	case c == '*':
		goto yystate40
	case c == '/':
		goto yystate43
	case c == '=':
		goto yystate44
	}

yystate40:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate41
	case c >= '\x01' && c <= ')' || c >= '+' && c <= 'ÿ':
		goto yystate40
	}

yystate41:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate41
	case c == '/':
		goto yystate42
	case c >= '\x01' && c <= ')' || c >= '+' && c <= '.' || c >= '0' && c <= 'ÿ':
		goto yystate40
	}

yystate42:
	c = s.next()
	goto yyrule3

yystate43:
	c = s.next()
	switch {
	default:
		goto yyrule4
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate43
	}

yystate44:
	c = s.next()
	goto yyrule28

yystate45:
	c = s.next()
	switch {
	default:
		goto yyrule77
	case c == '.':
		goto yystate34
	case c == '8' || c == '9':
		goto yystate47
	case c == 'E' || c == 'e':
		goto yystate35
	case c == 'X' || c == 'x':
		goto yystate49
	case c == 'i':
		goto yystate48
	case c >= '0' && c <= '7':
		goto yystate46
	}

yystate46:
	c = s.next()
	switch {
	default:
		goto yyrule77
	case c == '.':
		goto yystate34
	case c == '8' || c == '9':
		goto yystate47
	case c == 'E' || c == 'e':
		goto yystate35
	case c == 'i':
		goto yystate48
	case c >= '0' && c <= '7':
		goto yystate46
	}

yystate47:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate34
	case c == 'E' || c == 'e':
		goto yystate35
	case c == 'i':
		goto yystate48
	case c >= '0' && c <= '9':
		goto yystate47
	}

yystate48:
	c = s.next()
	goto yyrule78

yystate49:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate50
	}

yystate50:
	c = s.next()
	switch {
	default:
		goto yyrule77
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate50
	}

yystate51:
	c = s.next()
	switch {
	default:
		goto yyrule77
	case c == '.':
		goto yystate34
	case c == 'E' || c == 'e':
		goto yystate35
	case c == 'i':
		goto yystate48
	case c >= '0' && c <= '9':
		goto yystate51
	}

yystate52:
	c = s.next()
	switch {
	default:
		goto yyrule29
	case c == '=':
		goto yystate53
	}

yystate53:
	c = s.next()
	goto yyrule30

yystate54:
	c = s.next()
	goto yyrule31

yystate55:
	c = s.next()
	switch {
	default:
		goto yyrule32
	case c == '-':
		goto yystate56
	case c == '<':
		goto yystate57
	case c == '=':
		goto yystate59
	}

yystate56:
	c = s.next()
	goto yyrule33

yystate57:
	c = s.next()
	switch {
	default:
		goto yyrule34
	case c == '=':
		goto yystate58
	}

yystate58:
	c = s.next()
	goto yyrule35

yystate59:
	c = s.next()
	goto yyrule36

yystate60:
	c = s.next()
	switch {
	default:
		goto yyrule37
	case c == '=':
		goto yystate61
	}

yystate61:
	c = s.next()
	goto yyrule38

yystate62:
	c = s.next()
	switch {
	default:
		goto yyrule39
	case c == '=':
		goto yystate63
	case c == '>':
		goto yystate64
	}

yystate63:
	c = s.next()
	goto yyrule40

yystate64:
	c = s.next()
	switch {
	default:
		goto yyrule41
	case c == '=':
		goto yystate65
	}

yystate65:
	c = s.next()
	goto yyrule42

yystate66:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate67:
	c = s.next()
	goto yyrule43

yystate68:
	c = s.next()
	goto yyrule44

yystate69:
	c = s.next()
	switch {
	default:
		goto yyrule45
	case c == '=':
		goto yystate70
	}

yystate70:
	c = s.next()
	goto yyrule46

yystate71:
	c = s.next()
	goto yyrule82

yystate72:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'r':
		goto yystate73
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	}

yystate73:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate74
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate74:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate75
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	}

yystate75:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'k':
		goto yystate76
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	}

yystate76:
	c = s.next()
	switch {
	default:
		goto yyrule52
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate77:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate78
	case c == 'h':
		goto yystate81
	case c == 'o':
		goto yystate84
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'g' || c >= 'i' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	}

yystate78:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 's':
		goto yystate79
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	}

yystate79:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate80
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate80:
	c = s.next()
	switch {
	default:
		goto yyrule53
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate81:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate82
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	}

yystate82:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'n':
		goto yystate83
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	}

yystate83:
	c = s.next()
	switch {
	default:
		goto yyrule54
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate84:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'n':
		goto yystate85
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	}

yystate85:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 's':
		goto yystate86
	case c == 't':
		goto yystate88
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate86:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 't':
		goto yystate87
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate87:
	c = s.next()
	switch {
	default:
		goto yyrule55
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate88:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'i':
		goto yystate89
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	}

yystate89:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'n':
		goto yystate90
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	}

yystate90:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'u':
		goto yystate91
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	}

yystate91:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate92
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate92:
	c = s.next()
	switch {
	default:
		goto yyrule56
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate93:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate94
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate94:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'f':
		goto yystate95
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate66
	}

yystate95:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate96
	case c == 'e':
		goto yystate100
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate96:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'u':
		goto yystate97
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	}

yystate97:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'l':
		goto yystate98
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	}

yystate98:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 't':
		goto yystate99
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate99:
	c = s.next()
	switch {
	default:
		goto yyrule57
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate100:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'r':
		goto yystate101
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	}

yystate101:
	c = s.next()
	switch {
	default:
		goto yyrule58
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate102:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'l':
		goto yystate103
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	}

yystate103:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 's':
		goto yystate104
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate66
	}

yystate104:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate105
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate105:
	c = s.next()
	switch {
	default:
		goto yyrule59
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate106:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate107
	case c == 'o':
		goto yystate117
	case c == 'u':
		goto yystate119
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'n' || c >= 'p' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	}

yystate107:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'l':
		goto yystate108
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	}

yystate108:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'l':
		goto yystate109
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	}

yystate109:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 't':
		goto yystate110
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate110:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'h':
		goto yystate111
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	}

yystate111:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'r':
		goto yystate112
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	}

yystate112:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'o':
		goto yystate113
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	}

yystate113:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'u':
		goto yystate114
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	}

yystate114:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'g':
		goto yystate115
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	}

yystate115:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'h':
		goto yystate116
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	}

yystate116:
	c = s.next()
	switch {
	default:
		goto yyrule60
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate117:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'r':
		goto yystate118
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	}

yystate118:
	c = s.next()
	switch {
	default:
		goto yyrule61
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate119:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'n':
		goto yystate120
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	}

yystate120:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'c':
		goto yystate121
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	}

yystate121:
	c = s.next()
	switch {
	default:
		goto yyrule62
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate122:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'o':
		goto yystate123
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	}

yystate123:
	c = s.next()
	switch {
	default:
		goto yyrule63
	case c == 't':
		goto yystate124
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate124:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'o':
		goto yystate125
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	}

yystate125:
	c = s.next()
	switch {
	default:
		goto yyrule64
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate126:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'f':
		goto yystate127
	case c == 'm':
		goto yystate128
	case c == 'n':
		goto yystate133
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'l' || c >= 'o' && c <= 'z':
		goto yystate66
	}

yystate127:
	c = s.next()
	switch {
	default:
		goto yyrule65
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate128:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'p':
		goto yystate129
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	}

yystate129:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'o':
		goto yystate130
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate66
	}

yystate130:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'r':
		goto yystate131
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	}

yystate131:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 't':
		goto yystate132
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate132:
	c = s.next()
	switch {
	default:
		goto yyrule66
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate133:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 't':
		goto yystate134
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate134:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate135
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate135:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'r':
		goto yystate136
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	}

yystate136:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'f':
		goto yystate137
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate66
	}

yystate137:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate138
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	}

yystate138:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'c':
		goto yystate139
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	}

yystate139:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate140
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate140:
	c = s.next()
	switch {
	default:
		goto yyrule67
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate141:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate142
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	}

yystate142:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'p':
		goto yystate143
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	}

yystate143:
	c = s.next()
	switch {
	default:
		goto yyrule68
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate144:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate145
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	}

yystate145:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'c':
		goto yystate146
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	}

yystate146:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'k':
		goto yystate147
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate66
	}

yystate147:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate148
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	}

yystate148:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'g':
		goto yystate149
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	}

yystate149:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate150
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate150:
	c = s.next()
	switch {
	default:
		goto yyrule69
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate151:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate152
	case c == 'e':
		goto yystate156
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate152:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'n':
		goto yystate153
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	}

yystate153:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'g':
		goto yystate154
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate66
	}

yystate154:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate155
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate155:
	c = s.next()
	switch {
	default:
		goto yyrule70
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate156:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 't':
		goto yystate157
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate157:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'u':
		goto yystate158
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	}

yystate158:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'r':
		goto yystate159
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	}

yystate159:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'n':
		goto yystate160
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate66
	}

yystate160:
	c = s.next()
	switch {
	default:
		goto yyrule71
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate161:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate162
	case c == 't':
		goto yystate167
	case c == 'w':
		goto yystate172
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 's' || c == 'u' || c == 'v' || c >= 'x' && c <= 'z':
		goto yystate66
	}

yystate162:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'l':
		goto yystate163
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate66
	}

yystate163:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate164
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate164:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'c':
		goto yystate165
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	}

yystate165:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 't':
		goto yystate166
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate166:
	c = s.next()
	switch {
	default:
		goto yyrule72
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate167:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'r':
		goto yystate168
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	}

yystate168:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'u':
		goto yystate169
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate66
	}

yystate169:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'c':
		goto yystate170
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	}

yystate170:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 't':
		goto yystate171
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate171:
	c = s.next()
	switch {
	default:
		goto yyrule73
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate172:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'i':
		goto yystate173
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate66
	}

yystate173:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 't':
		goto yystate174
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate66
	}

yystate174:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'c':
		goto yystate175
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate66
	}

yystate175:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'h':
		goto yystate176
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate66
	}

yystate176:
	c = s.next()
	switch {
	default:
		goto yyrule74
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate177:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'y':
		goto yystate178
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate66
	}

yystate178:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'p':
		goto yystate179
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate66
	}

yystate179:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'e':
		goto yystate180
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate66
	}

yystate180:
	c = s.next()
	switch {
	default:
		goto yyrule75
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate181:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'a':
		goto yystate182
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate66
	}

yystate182:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == 'r':
		goto yystate183
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate66
	}

yystate183:
	c = s.next()
	switch {
	default:
		goto yyrule76
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate66
	}

yystate184:
	c = s.next()
	goto yyrule47

yystate185:
	c = s.next()
	switch {
	default:
		goto yyrule48
	case c == '=':
		goto yystate186
	case c == '|':
		goto yystate187
	}

yystate186:
	c = s.next()
	goto yyrule49

yystate187:
	c = s.next()
	goto yyrule50

yystate188:
	c = s.next()
	goto yyrule51

yystate189:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate190
	}

yystate190:
	c = s.next()
	goto yyrule89

yystate191:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate189
	}

yystate192:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate191
	}

	goto yystate193 // silence unused label error
yystate193:
	c = s.next()
yystart193:
	switch {
	default:
		goto yystate194 // c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ'
	case c == '"':
		goto yystate195
	case c == '\\':
		goto yystate196
	case c == '\x00':
		goto yystate2
	}

yystate194:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate195
	case c == '\\':
		goto yystate196
	case c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate194
	}

yystate195:
	c = s.next()
	goto yyrule86

yystate196:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate194
	}

	goto yystate197 // silence unused label error
yystate197:
	c = s.next()
yystart197:
	switch {
	default:
		goto yystate198 // c >= '\x01' && c <= '_' || c >= 'a' && c <= 'ÿ'
	case c == '\x00':
		goto yystate2
	case c == '`':
		goto yystate199
	}

yystate198:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '`':
		goto yystate199
	case c >= '\x01' && c <= '_' || c >= 'a' && c <= 'ÿ':
		goto yystate198
	}

yystate199:
	c = s.next()
	goto yyrule87

yyrule1: // \0
	{
		s.i0++
		return token.EOF, lval
	}
yyrule2: // [ \t\r]+

	goto yystate0
yyrule3: // \/\*([^*]|\*+[^*/])*\*+\/
yyrule4: // \/\/.*
	{
		return token.COMMENT, string(s.val)
	}
yyrule5: // "!"
	{
		return token.NOT, lval
	}
yyrule6: // "!="
	{
		return token.NEQ, lval
	}
yyrule7: // "%"
	{
		return token.REM, lval
	}
yyrule8: // "%="
	{
		return token.REM_ASSIGN, lval
	}
yyrule9: // "&"
	{
		return token.AND, lval
	}
yyrule10: // "&&"
	{
		return token.LAND, lval
	}
yyrule11: // "&="
	{
		return token.AND_ASSIGN, lval
	}
yyrule12: // "&^"
	{
		return token.AND_NOT, lval
	}
yyrule13: // "&^="
	{
		return token.AND_NOT_ASSIGN, lval
	}
yyrule14: // "("
	{
		return token.LPAREN, lval
	}
yyrule15: // ")"
	{
		return token.RPAREN, lval
	}
yyrule16: // "*"
	{
		return token.MUL, lval
	}
yyrule17: // "*="
	{
		return token.MUL_ASSIGN, lval
	}
yyrule18: // "+"
	{
		return token.ADD, lval
	}
yyrule19: // "++"
	{
		return token.INC, lval
	}
yyrule20: // "+="
	{
		return token.ADD_ASSIGN, lval
	}
yyrule21: // ","
	{
		return token.COMMA, lval
	}
yyrule22: // "-"
	{
		return token.SUB, lval
	}
yyrule23: // "--"
	{
		return token.DEC, lval
	}
yyrule24: // "-="
	{
		return token.SUB_ASSIGN, lval
	}
yyrule25: // "."
	{
		return token.PERIOD, lval
	}
yyrule26: // "..."
	{
		return token.ELLIPSIS, lval
	}
yyrule27: // "/"
	{
		return token.QUO, lval
	}
yyrule28: // "/="
	{
		return token.QUO_ASSIGN, lval
	}
yyrule29: // ":"
	{
		return token.COLON, lval
	}
yyrule30: // ":="
	{
		return token.DEFINE, lval
	}
yyrule31: // ";"
	{
		return token.SEMICOLON, lval
	}
yyrule32: // "<"
	{
		return token.LSS, lval
	}
yyrule33: // "<-"
	{
		return token.ARROW, lval
	}
yyrule34: // "<<"
	{
		return token.SHL, lval
	}
yyrule35: // "<<="
	{
		return token.SHL_ASSIGN, lval
	}
yyrule36: // "<="
	{
		return token.LEQ, lval
	}
yyrule37: // "="
	{
		return token.ASSIGN, lval
	}
yyrule38: // "=="
	{
		return token.EQL, lval
	}
yyrule39: // ">"
	{
		return token.GTR, lval
	}
yyrule40: // ">="
	{
		return token.GEQ, lval
	}
yyrule41: // ">>"
	{
		return token.SHR, lval
	}
yyrule42: // ">>="
	{
		return token.SHR_ASSIGN, lval
	}
yyrule43: // "["
	{
		return token.LBRACK, lval
	}
yyrule44: // "]"
	{
		return token.RBRACK, lval
	}
yyrule45: // "^"
	{
		return token.XOR, lval
	}
yyrule46: // "^="
	{
		return token.XOR_ASSIGN, lval
	}
yyrule47: // "{"
	{
		return token.LBRACE, lval
	}
yyrule48: // "|"
	{
		return token.OR, lval
	}
yyrule49: // "|="
	{
		return token.OR_ASSIGN, lval
	}
yyrule50: // "||"
	{
		return token.LOR, lval
	}
yyrule51: // "}"
	{
		return token.RBRACE, lval
	}
yyrule52: // break
	{
		return token.BREAK, lval
	}
yyrule53: // case
	{
		return token.CASE, lval
	}
yyrule54: // chan
	{
		return token.CHAN, lval
	}
yyrule55: // const
	{
		return token.CONST, lval
	}
yyrule56: // continue
	{
		return token.CONTINUE, lval
	}
yyrule57: // default
	{
		return token.DEFAULT, lval
	}
yyrule58: // defer
	{
		return token.DEFER, lval
	}
yyrule59: // else
	{
		return token.ELSE, lval
	}
yyrule60: // fallthrough
	{
		return token.FALLTHROUGH, lval
	}
yyrule61: // for
	{
		return token.FOR, lval
	}
yyrule62: // func
	{
		return token.FUNC, lval
	}
yyrule63: // go
	{
		return token.GO, lval
	}
yyrule64: // goto
	{
		return token.GOTO, lval
	}
yyrule65: // if
	{
		return token.IF, lval
	}
yyrule66: // import
	{
		return token.IMPORT, lval
	}
yyrule67: // interface
	{
		return token.INTERFACE, lval
	}
yyrule68: // map
	{
		return token.MAP, lval
	}
yyrule69: // package
	{
		return token.PACKAGE, lval
	}
yyrule70: // range
	{
		return token.RANGE, lval
	}
yyrule71: // return
	{
		return token.RETURN, lval
	}
yyrule72: // select
	{
		return token.SELECT, lval
	}
yyrule73: // struct
	{
		return token.STRUCT, lval
	}
yyrule74: // switch
	{
		return token.SWITCH, lval
	}
yyrule75: // type
	{
		return token.TYPE, lval
	}
yyrule76: // var
	{
		return token.VAR, lval
	}
yyrule77: // {int_lit}
	{
		return token.INT, string(s.val)
	}
yyrule78: // {imaginary_ilit}
yyrule79: // {imaginary_lit}
	{
		return token.IMAG, string(s.val)
	}
yyrule80: // {float_lit}
	{
		return token.FLOAT, string(s.val)
	}
yyrule81: // \"
	{
		s.sc = S1
		goto yystate0
	}
yyrule82: // `
	{
		s.sc = S2
		goto yystate0
	}
yyrule83: // ''
yyrule84: // '(\\.)?[^']*
yyrule85: // '(\\.)?[^']*'
	{
		return token.CHAR, string(s.val)
	}
yyrule86: // (\\.|[^\\"])*\"
	{
		s.Col--
		s.i0--
		return token.STRING, `"` + string(s.val)
	}
yyrule87: // ([^`]|\n)*`
	{
		s.Col--
		s.i0--
		return token.STRING, "`" + string(s.val)
	}
yyrule88: // [a-zA-Z_][a-zA-Z_0-9]*
	{

		if c >= '\xC2' && c <= '\xF4' {
			s.i--
			s.NCol--
			for rune := rune(1); rune >= 0; rune = s.getRune(true) {
			}
			tok, lval = token.IDENT, string(s.src[s.i0-1:s.i])
			s.next()
			return
		}
		return token.IDENT, string(s.val)
	}
yyrule89: // {non_ascii}
	{

		s.i = s.i0 - 1
		if rune := s.getRune(false); rune < 0 {
			_, sz := utf8.DecodeRune(s.src[s.i:])
			s.i += sz
			s.next()
			s.err("expected unicode lettter, got %U", rune)
			return token.ILLEGAL, string(-rune)
		}
		for rune := rune(1); rune >= 0; rune = s.getRune(true) {
		}
		s.next()
		return token.IDENT, string(s.src[s.i0-1 : s.i-1])
	}
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	s.next()
	return token.ILLEGAL, string(c0)
}

func (s *Scanner) getRune(acceptDigits bool) (r rune) {
	var sz int
	if r, sz = utf8.DecodeRune(s.src[s.i:]); sz != 0 &&
		(r == '_' || unicode.IsLetter(r) || (acceptDigits && unicode.IsDigit(r))) {
		s.i += sz
		s.NCol += sz
		return
	}

	return -r
}
