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
	"fmt"
	"go/token"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// A Scanner holds the scanner's internal state while processing a given text.
type Scanner struct {
	Col    int     // Starting column of the last scanned token.
	Errors []error // List of accumulated errors.
	Fname  string  // File name (reported) of the scanned source.
	Line   int     // Starting line of the last scanned token.
	NCol   int     // Starting column (reported) for the next scanned token.
	NLine  int     // Starting line (reported) for the next scanned token.
	c      int
	i      int
	i0     int
	sc     int
	src    []byte
	val    []byte
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

// Error appends s.Fname:s.Line:s.Col msg to s.Errors.
func (s *Scanner) Error(msg string) {
	s.err(msg)
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
		goto yystart192
	case 2: // start condition: S2
		goto yystart196
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
		goto yystate19
	case c == ')':
		goto yystate20
	case c == '*':
		goto yystate21
	case c == '+':
		goto yystate23
	case c == ',':
		goto yystate26
	case c == '-':
		goto yystate27
	case c == '.':
		goto yystate30
	case c == '/':
		goto yystate38
	case c == '0':
		goto yystate44
	case c == ':':
		goto yystate51
	case c == ';':
		goto yystate53
	case c == '<':
		goto yystate54
	case c == '=':
		goto yystate59
	case c == '>':
		goto yystate61
	case c == '[':
		goto yystate66
	case c == '\'':
		goto yystate14
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate3
	case c == '\x00':
		goto yystate2
	case c == ']':
		goto yystate67
	case c == '^':
		goto yystate68
	case c == '`':
		goto yystate70
	case c == 'b':
		goto yystate71
	case c == 'c':
		goto yystate76
	case c == 'd':
		goto yystate92
	case c == 'e':
		goto yystate101
	case c == 'f':
		goto yystate105
	case c == 'g':
		goto yystate121
	case c == 'i':
		goto yystate125
	case c == 'm':
		goto yystate140
	case c == 'p':
		goto yystate143
	case c == 'r':
		goto yystate150
	case c == 's':
		goto yystate160
	case c == 't':
		goto yystate176
	case c == 'v':
		goto yystate180
	case c == '{':
		goto yystate183
	case c == '|':
		goto yystate184
	case c == '}':
		goto yystate187
	case c >= '1' && c <= '9':
		goto yystate50
	case c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'h' || c >= 'j' && c <= 'l' || c == 'n' || c == 'o' || c == 'q' || c == 'u' || c >= 'w' && c <= 'z':
		goto yystate65
	case c >= 'Â' && c <= 'ß':
		goto yystate188
	case c >= 'à' && c <= 'ï':
		goto yystate190
	case c >= 'ð' && c <= 'ô':
		goto yystate191
	}

yystate2:
	c = s.next()
	goto yyrule1

yystate3:
	c = s.next()
	switch {
	default:
		goto yyrule2
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
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
		goto yyabort
	case c == '\'':
		goto yystate16
	case c == '\\':
		goto yystate17
	case c >= '\x01' && c <= '&' || c >= '(' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate15
	}

yystate15:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '\'':
		goto yystate16
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate15
	}

yystate16:
	c = s.next()
	goto yyrule83

yystate17:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '\'':
		goto yystate18
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate15
	}

yystate18:
	c = s.next()
	switch {
	default:
		goto yyrule83
	case c == '\'':
		goto yystate16
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate15
	}

yystate19:
	c = s.next()
	goto yyrule14

yystate20:
	c = s.next()
	goto yyrule15

yystate21:
	c = s.next()
	switch {
	default:
		goto yyrule16
	case c == '=':
		goto yystate22
	}

yystate22:
	c = s.next()
	goto yyrule17

yystate23:
	c = s.next()
	switch {
	default:
		goto yyrule18
	case c == '+':
		goto yystate24
	case c == '=':
		goto yystate25
	}

yystate24:
	c = s.next()
	goto yyrule19

yystate25:
	c = s.next()
	goto yyrule20

yystate26:
	c = s.next()
	goto yyrule21

yystate27:
	c = s.next()
	switch {
	default:
		goto yyrule22
	case c == '-':
		goto yystate28
	case c == '=':
		goto yystate29
	}

yystate28:
	c = s.next()
	goto yyrule23

yystate29:
	c = s.next()
	goto yyrule24

yystate30:
	c = s.next()
	switch {
	default:
		goto yyrule25
	case c == '.':
		goto yystate31
	case c >= '0' && c <= '9':
		goto yystate33
	}

yystate31:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate32
	}

yystate32:
	c = s.next()
	goto yyrule26

yystate33:
	c = s.next()
	switch {
	default:
		goto yyrule80
	case c == 'E' || c == 'e':
		goto yystate34
	case c == 'i':
		goto yystate37
	case c >= '0' && c <= '9':
		goto yystate33
	}

yystate34:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '+' || c == '-':
		goto yystate35
	case c >= '0' && c <= '9':
		goto yystate36
	}

yystate35:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate36
	}

yystate36:
	c = s.next()
	switch {
	default:
		goto yyrule80
	case c == 'i':
		goto yystate37
	case c >= '0' && c <= '9':
		goto yystate36
	}

yystate37:
	c = s.next()
	goto yyrule78

yystate38:
	c = s.next()
	switch {
	default:
		goto yyrule27
	case c == '*':
		goto yystate39
	case c == '/':
		goto yystate42
	case c == '=':
		goto yystate43
	}

yystate39:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate40
	case c >= '\x01' && c <= ')' || c >= '+' && c <= 'ÿ':
		goto yystate39
	}

yystate40:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate40
	case c == '/':
		goto yystate41
	case c >= '\x01' && c <= ')' || c >= '+' && c <= '.' || c >= '0' && c <= 'ÿ':
		goto yystate39
	}

yystate41:
	c = s.next()
	goto yyrule3

yystate42:
	c = s.next()
	switch {
	default:
		goto yyrule4
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate42
	}

yystate43:
	c = s.next()
	goto yyrule28

yystate44:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c == '.':
		goto yystate33
	case c == '8' || c == '9':
		goto yystate46
	case c == 'E' || c == 'e':
		goto yystate34
	case c == 'X' || c == 'x':
		goto yystate48
	case c == 'i':
		goto yystate47
	case c >= '0' && c <= '7':
		goto yystate45
	}

yystate45:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c == '.':
		goto yystate33
	case c == '8' || c == '9':
		goto yystate46
	case c == 'E' || c == 'e':
		goto yystate34
	case c == 'i':
		goto yystate47
	case c >= '0' && c <= '7':
		goto yystate45
	}

yystate46:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate33
	case c == 'E' || c == 'e':
		goto yystate34
	case c == 'i':
		goto yystate47
	case c >= '0' && c <= '9':
		goto yystate46
	}

yystate47:
	c = s.next()
	goto yyrule77

yystate48:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate49
	}

yystate49:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate49
	}

yystate50:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c == '.':
		goto yystate33
	case c == 'E' || c == 'e':
		goto yystate34
	case c == 'i':
		goto yystate47
	case c >= '0' && c <= '9':
		goto yystate50
	}

yystate51:
	c = s.next()
	switch {
	default:
		goto yyrule29
	case c == '=':
		goto yystate52
	}

yystate52:
	c = s.next()
	goto yyrule30

yystate53:
	c = s.next()
	goto yyrule31

yystate54:
	c = s.next()
	switch {
	default:
		goto yyrule32
	case c == '-':
		goto yystate55
	case c == '<':
		goto yystate56
	case c == '=':
		goto yystate58
	}

yystate55:
	c = s.next()
	goto yyrule33

yystate56:
	c = s.next()
	switch {
	default:
		goto yyrule34
	case c == '=':
		goto yystate57
	}

yystate57:
	c = s.next()
	goto yyrule35

yystate58:
	c = s.next()
	goto yyrule36

yystate59:
	c = s.next()
	switch {
	default:
		goto yyrule37
	case c == '=':
		goto yystate60
	}

yystate60:
	c = s.next()
	goto yyrule38

yystate61:
	c = s.next()
	switch {
	default:
		goto yyrule39
	case c == '=':
		goto yystate62
	case c == '>':
		goto yystate63
	}

yystate62:
	c = s.next()
	goto yyrule40

yystate63:
	c = s.next()
	switch {
	default:
		goto yyrule41
	case c == '=':
		goto yystate64
	}

yystate64:
	c = s.next()
	goto yyrule42

yystate65:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate66:
	c = s.next()
	goto yyrule43

yystate67:
	c = s.next()
	goto yyrule44

yystate68:
	c = s.next()
	switch {
	default:
		goto yyrule45
	case c == '=':
		goto yystate69
	}

yystate69:
	c = s.next()
	goto yyrule46

yystate70:
	c = s.next()
	goto yyrule82

yystate71:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'r':
		goto yystate72
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate65
	}

yystate72:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate73
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate73:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate74
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate65
	}

yystate74:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'k':
		goto yystate75
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate65
	}

yystate75:
	c = s.next()
	switch {
	default:
		goto yyrule52
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate76:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate77
	case c == 'h':
		goto yystate80
	case c == 'o':
		goto yystate83
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'g' || c >= 'i' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate65
	}

yystate77:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 's':
		goto yystate78
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate65
	}

yystate78:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate79
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate79:
	c = s.next()
	switch {
	default:
		goto yyrule53
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate80:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate81
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate65
	}

yystate81:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'n':
		goto yystate82
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate65
	}

yystate82:
	c = s.next()
	switch {
	default:
		goto yyrule54
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate83:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'n':
		goto yystate84
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate65
	}

yystate84:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 's':
		goto yystate85
	case c == 't':
		goto yystate87
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate85:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 't':
		goto yystate86
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate86:
	c = s.next()
	switch {
	default:
		goto yyrule55
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate87:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'i':
		goto yystate88
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate65
	}

yystate88:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'n':
		goto yystate89
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate65
	}

yystate89:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'u':
		goto yystate90
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate65
	}

yystate90:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate91
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate91:
	c = s.next()
	switch {
	default:
		goto yyrule56
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate92:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate93
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate93:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'f':
		goto yystate94
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate65
	}

yystate94:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate95
	case c == 'e':
		goto yystate99
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate95:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'u':
		goto yystate96
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate65
	}

yystate96:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'l':
		goto yystate97
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate65
	}

yystate97:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 't':
		goto yystate98
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate98:
	c = s.next()
	switch {
	default:
		goto yyrule57
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate99:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'r':
		goto yystate100
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate65
	}

yystate100:
	c = s.next()
	switch {
	default:
		goto yyrule58
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate101:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'l':
		goto yystate102
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate65
	}

yystate102:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 's':
		goto yystate103
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate65
	}

yystate103:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate104
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate104:
	c = s.next()
	switch {
	default:
		goto yyrule59
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate105:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate106
	case c == 'o':
		goto yystate116
	case c == 'u':
		goto yystate118
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'n' || c >= 'p' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate65
	}

yystate106:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'l':
		goto yystate107
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate65
	}

yystate107:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'l':
		goto yystate108
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate65
	}

yystate108:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 't':
		goto yystate109
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate109:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'h':
		goto yystate110
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate65
	}

yystate110:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'r':
		goto yystate111
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate65
	}

yystate111:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'o':
		goto yystate112
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate65
	}

yystate112:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'u':
		goto yystate113
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate65
	}

yystate113:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'g':
		goto yystate114
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate65
	}

yystate114:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'h':
		goto yystate115
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate65
	}

yystate115:
	c = s.next()
	switch {
	default:
		goto yyrule60
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate116:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'r':
		goto yystate117
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate65
	}

yystate117:
	c = s.next()
	switch {
	default:
		goto yyrule61
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate118:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'n':
		goto yystate119
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate65
	}

yystate119:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'c':
		goto yystate120
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate65
	}

yystate120:
	c = s.next()
	switch {
	default:
		goto yyrule62
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate121:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'o':
		goto yystate122
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate65
	}

yystate122:
	c = s.next()
	switch {
	default:
		goto yyrule63
	case c == 't':
		goto yystate123
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate123:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'o':
		goto yystate124
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate65
	}

yystate124:
	c = s.next()
	switch {
	default:
		goto yyrule64
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate125:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'f':
		goto yystate126
	case c == 'm':
		goto yystate127
	case c == 'n':
		goto yystate132
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'l' || c >= 'o' && c <= 'z':
		goto yystate65
	}

yystate126:
	c = s.next()
	switch {
	default:
		goto yyrule65
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate127:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'p':
		goto yystate128
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate65
	}

yystate128:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'o':
		goto yystate129
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate65
	}

yystate129:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'r':
		goto yystate130
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate65
	}

yystate130:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 't':
		goto yystate131
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate131:
	c = s.next()
	switch {
	default:
		goto yyrule66
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate132:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 't':
		goto yystate133
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate133:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate134
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate134:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'r':
		goto yystate135
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate65
	}

yystate135:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'f':
		goto yystate136
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate65
	}

yystate136:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate137
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate65
	}

yystate137:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'c':
		goto yystate138
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate65
	}

yystate138:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate139
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate139:
	c = s.next()
	switch {
	default:
		goto yyrule67
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate140:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate141
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate65
	}

yystate141:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'p':
		goto yystate142
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate65
	}

yystate142:
	c = s.next()
	switch {
	default:
		goto yyrule68
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate143:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate144
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate65
	}

yystate144:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'c':
		goto yystate145
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate65
	}

yystate145:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'k':
		goto yystate146
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate65
	}

yystate146:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate147
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate65
	}

yystate147:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'g':
		goto yystate148
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate65
	}

yystate148:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate149
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate149:
	c = s.next()
	switch {
	default:
		goto yyrule69
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate150:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate151
	case c == 'e':
		goto yystate155
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate151:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'n':
		goto yystate152
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate65
	}

yystate152:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'g':
		goto yystate153
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate65
	}

yystate153:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate154
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate154:
	c = s.next()
	switch {
	default:
		goto yyrule70
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate155:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 't':
		goto yystate156
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate156:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'u':
		goto yystate157
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate65
	}

yystate157:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'r':
		goto yystate158
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate65
	}

yystate158:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'n':
		goto yystate159
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate65
	}

yystate159:
	c = s.next()
	switch {
	default:
		goto yyrule71
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate160:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate161
	case c == 't':
		goto yystate166
	case c == 'w':
		goto yystate171
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 's' || c == 'u' || c == 'v' || c >= 'x' && c <= 'z':
		goto yystate65
	}

yystate161:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'l':
		goto yystate162
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate65
	}

yystate162:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate163
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate163:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'c':
		goto yystate164
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate65
	}

yystate164:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 't':
		goto yystate165
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate165:
	c = s.next()
	switch {
	default:
		goto yyrule72
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate166:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'r':
		goto yystate167
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate65
	}

yystate167:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'u':
		goto yystate168
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate65
	}

yystate168:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'c':
		goto yystate169
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate65
	}

yystate169:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 't':
		goto yystate170
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate170:
	c = s.next()
	switch {
	default:
		goto yyrule73
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate171:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'i':
		goto yystate172
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate65
	}

yystate172:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 't':
		goto yystate173
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate65
	}

yystate173:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'c':
		goto yystate174
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate65
	}

yystate174:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'h':
		goto yystate175
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate65
	}

yystate175:
	c = s.next()
	switch {
	default:
		goto yyrule74
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate176:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'y':
		goto yystate177
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate65
	}

yystate177:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'p':
		goto yystate178
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate65
	}

yystate178:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'e':
		goto yystate179
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate65
	}

yystate179:
	c = s.next()
	switch {
	default:
		goto yyrule75
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate180:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'a':
		goto yystate181
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate65
	}

yystate181:
	c = s.next()
	switch {
	default:
		goto yyrule86
	case c == 'r':
		goto yystate182
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate65
	}

yystate182:
	c = s.next()
	switch {
	default:
		goto yyrule76
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate65
	}

yystate183:
	c = s.next()
	goto yyrule47

yystate184:
	c = s.next()
	switch {
	default:
		goto yyrule48
	case c == '=':
		goto yystate185
	case c == '|':
		goto yystate186
	}

yystate185:
	c = s.next()
	goto yyrule49

yystate186:
	c = s.next()
	goto yyrule50

yystate187:
	c = s.next()
	goto yyrule51

yystate188:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate189
	}

yystate189:
	c = s.next()
	goto yyrule87

yystate190:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate188
	}

yystate191:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate190
	}

	goto yystate192 // silence unused label error
yystate192:
	c = s.next()
yystart192:
	switch {
	default:
		goto yystate193 // c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ'
	case c == '"':
		goto yystate194
	case c == '\\':
		goto yystate195
	case c == '\x00':
		goto yystate2
	}

yystate193:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate194
	case c == '\\':
		goto yystate195
	case c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate193
	}

yystate194:
	c = s.next()
	goto yyrule84

yystate195:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate193
	}

	goto yystate196 // silence unused label error
yystate196:
	c = s.next()
yystart196:
	switch {
	default:
		goto yystate197 // c >= '\x01' && c <= '_' || c >= 'a' && c <= 'ÿ'
	case c == '\x00':
		goto yystate2
	case c == '`':
		goto yystate198
	}

yystate197:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '`':
		goto yystate198
	case c >= '\x01' && c <= '_' || c >= 'a' && c <= 'ÿ':
		goto yystate197
	}

yystate198:
	c = s.next()
	goto yyrule85

yyrule1: // \0
	{
		s.i0++
		return token.EOF, lval
	}
yyrule2: // [ \t\n\r]+

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
yyrule77: // {imaginary_ilit}
	{
		return s.int(token.IMAG)
	}
yyrule78: // {imaginary_lit}
	{
		return s.float(token.IMAG)
	}
yyrule79: // {int_lit}
	{
		return s.int(token.INT)
	}
yyrule80: // {float_lit}
	{
		return s.float(token.FLOAT)
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
yyrule83: // '(\\.)?[^']*'
	{
		if tok, lval = s.str(""); tok != token.STRING {
			return
		}
		s.i0++
		return token.CHAR, int32(lval.(string)[0])
	}
yyrule84: // (\\.|[^\\"])*\"
	{
		return s.str("\"")
	}
yyrule85: // ([^`]|\n)*`
	{
		return s.str("`")
	}
yyrule86: // [a-zA-Z_][a-zA-Z_0-9]*
	{

		if c >= '\xC2' && c <= '\xF4' {
			s.i--
			s.NCol--
			for rune := rune(1); rune >= 0; rune = s.getRune(true) {
			}
			s.next()
		}
		if s.i < len(s.src) {
			return token.IDENT, string(s.src[s.i0-1 : s.i-1])
		}

		return token.IDENT, string(s.src[s.i0-1 : s.i])
	}
yyrule87: // {non_ascii}
	{

		s.i = s.i0 - 1
		if rune := s.getRune(false); rune < 0 {
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

func (s *Scanner) str(pref string) (tok token.Token, lval interface{}) {
	s.sc = 0
	ss := pref + string(s.val)
	ss, err := strconv.Unquote(ss)
	if err != nil {
		s.err("string literal %q: %v", ss, err)
		return token.ILLEGAL, ss
	}

	s.i0--
	return token.STRING, ss
}

func (s *Scanner) int(tk token.Token) (tok token.Token, lval interface{}) {
	tok = tk
	if tok == token.IMAG {
		s.val = s.val[:len(s.val)-1]
	}
	n, err := strconv.ParseUint(string(s.val), 0, 64)
	if err != nil {
		lval = string(s.val)
	} else if tok == token.IMAG {
		lval = complex(0, float64(n))
	} else {
		lval = n
	}
	return
}

func (s *Scanner) float(tk token.Token) (tok token.Token, lval interface{}) {
	tok = tk
	if tok == token.IMAG {
		s.val = s.val[:len(s.val)-1]
	}
	n, err := strconv.ParseFloat(string(s.val), 64)
	if err != nil {
		lval = string(s.val)
	} else if tok == token.IMAG {
		lval = complex(0, n)
	} else {
		lval = n
	}
	return
}
