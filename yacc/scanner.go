// Copyright (c) 2013 Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// CAUTION: If this file is a Go source file (*.go), it was generated
// automatically by '$ golex' from a *.l file - DO NOT EDIT in that case!

// Package scanner (WIP:TODO) implements a scanner for yacc[1] source text with
// actions written in Go. It takes a []byte as source which can then be
// tokenized through repeated calls to the Scan method.
//
// Links
//
// Referenced from above:
//
// [1]: http://pubs.opengroup.org/onlinepubs/009695399/utilities/yacc.html
package scanner

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Token int

const (
	_INITIAL = iota
	_GO
	_S1
	_S2
)

const (
	_ = 0xE000 + iota

	// ------------------------------------------- Go mode or shared tokens

	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"

	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	// Keywords
	BREAK
	CASE
	CHAN
	CONST
	CONTINUE

	DEFAULT
	DEFER
	ELSE
	FALLTHROUGH
	FOR

	FUNC
	GO
	GOTO
	IF
	IMPORT

	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN

	SELECT
	STRUCT
	SWITCH
	GO_TYPE
	VAR

	// yacc mode tokens
	C_IDENTIFIER // IDENTIFIER ':'
	ERR_VERBOSE  // %error-verbose
	IDENTIFIER   // [a-zA-Z_][a-zA-Z0-9_.]*
	LCURL        // %{
	LEFT         // %left
	MARK         // %%
	NONASSOC     // %nonassoc
	PREC         // %prec
	RCURL        // %}
	RIGHT        // %right
	START        // %start
	TOKEN        // %token
	TYPE         // %type
	UNION        // %union
)

var ts = map[Token]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "ADD",
	SUB: "SUB",
	MUL: "MUL",
	QUO: "QUO",
	REM: "REM",

	AND:     "AND",
	OR:      "OR",
	XOR:     "XOR",
	SHL:     "SHL",
	SHR:     "SHR",
	AND_NOT: "AND_NOT",

	ADD_ASSIGN: "ADD_ASSIGN",
	SUB_ASSIGN: "SUB_ASSIGN",
	MUL_ASSIGN: "MUL_ASSIGN",
	QUO_ASSIGN: "QUO_ASSIGN",
	REM_ASSIGN: "REM_ASSIGN",

	AND_ASSIGN:     "AND_ASSIGN",
	OR_ASSIGN:      "OR_ASSIGN",
	XOR_ASSIGN:     "XOR_ASSIGN",
	SHL_ASSIGN:     "SHL_ASSIGN",
	SHR_ASSIGN:     "SHR_ASSIGN",
	AND_NOT_ASSIGN: "AND_NOT_ASSIGN",

	LAND:  "LAND",
	LOR:   "LOR",
	ARROW: "ARROW",
	INC:   "INC",
	DEC:   "DEC",

	EQL:    "EQL",
	LSS:    "LSS",
	GTR:    "GTR",
	ASSIGN: "ASSIGN",
	NOT:    "NOT",

	NEQ:      "NEQ",
	LEQ:      "LEQ",
	GEQ:      "GEQ",
	DEFINE:   "DEFINE",
	ELLIPSIS: "ELLIPSIS",

	LPAREN: "LPAREN",
	LBRACK: "LBRACK",
	LBRACE: "LBRACE",
	COMMA:  "COMMA",
	PERIOD: "PERIOD",

	RPAREN:    "RPAREN",
	RBRACK:    "RBRACK",
	RBRACE:    "RBRACE",
	SEMICOLON: "SEMICOLON",
	COLON:     "COLON",

	BREAK:    "BREAK",
	CASE:     "CASE",
	CHAN:     "CHAN",
	CONST:    "CONST",
	CONTINUE: "CONTINUE",

	DEFAULT:     "DEFAULT",
	DEFER:       "DEFER",
	ELSE:        "ELSE",
	FALLTHROUGH: "FALLTHROUGH",
	FOR:         "FOR",

	FUNC:   "FUNC",
	GO:     "GO",
	GOTO:   "GOTO",
	IF:     "IF",
	IMPORT: "IMPORT",

	INTERFACE: "INTERFACE",
	MAP:       "MAP",
	PACKAGE:   "PACKAGE",
	RANGE:     "RANGE",
	RETURN:    "RETURN",

	SELECT:  "SELECT",
	STRUCT:  "STRUCT",
	SWITCH:  "SWITCH",
	GO_TYPE: "GO_TYPE",
	VAR:     "VAR",

	// --------------------------------------------- yacc/bison mode tokens

	C_IDENTIFIER: "C_IDENTIFIER",
	IDENTIFIER:   "IDENTIFIER",
	LCURL:        "LCURL",
	LEFT:         "LEFT",
	MARK:         "MARK",

	NONASSOC: "NONASSOC",
	PREC:     "PREC",
	RCURL:    "RCURL",
	RIGHT:    "RIGHT",

	START:       "START",
	TOKEN:       "TOKEN",
	TYPE:        "TYPE",
	UNION:       "UNION",
	ERR_VERBOSE: "ERR_VERBOSE",
}

// String implements fmt.Stringer
func (i Token) String() string {
	if s := ts[i]; s != "" {
		return s
	}

	return fmt.Sprintf("%T(%d)", i, i)
}

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
	ssc    int // saved state condition
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

// Error implements yyLexer.
func (s *Scanner) Error(msg string) {
	switch msg {
	case "syntax error":
		s.err(msg)
	default:
		s.Errors = append(s.Errors, errors.New(msg))
	}
}

// Mode allows to switch the scanner mode from scanning yacc tokens to scanning
// rule action tokens (Go tokens). Initially the yacc mode is active.
func (s *Scanner) Mode(yacc bool) {
	switch yacc {
	case true:
		s.sc = _INITIAL
	case false:
		s.sc = _GO
	}
}

// Scan works like ScanRaw but recognizes also yacc's C_IDENTIFIER token (in
// yacc mode).
func (s *Scanner) Scan() (tok Token, lval interface{}) {
	tok, lval = s.ScanRaw()
	if s.sc != _INITIAL || tok != IDENTIFIER {
		return
	}

	i, nl, nc, c := s.i, s.NLine, s.NCol, s.c
	tok2, lit := s.ScanRaw()
	if tok2 != ILLEGAL || lit.(string) != ":" {
		s.i, s.NLine, s.NCol, s.c = i, nl, nc, c
		return
	}

	return C_IDENTIFIER, lval
}

// ScanRaw scans the next token and returns the token and its value if
// applicable.  The source end is indicated by EOF.
//
// If the returned token is IDENT, INT, FLOAT, IMAG, CHAR, STRING or COMMENT,
// lval has has the corresponding value - not the string representation of the
// value. However, numeric literals which overflow the corresponding Go
// predeclared types are returned as string.
//
// If the returned token is ILLEGAL, the literal string is the offending
// character or number/string/char literal.
func (s *Scanner) ScanRaw() (tok Token, lval interface{}) {
	//defer func() { fmt.Printf("%s(%d) %v\n", tok, int(tok), lval) }()
	c0, c := s.c, s.c

yystate0:

	s.val = s.val[:0]
	s.i0, s.Line, s.Col, c0 = s.i, s.NLine, s.NCol, c

	switch yyt := s.sc; yyt {
	default:
		panic(fmt.Errorf(`invalid start condition %d`, yyt))
	case 0: // start condition: INITIAL
		goto yystart1
	case 1: // start condition: _GO
		goto yystart83
	case 2: // start condition: _S1
		goto yystart259
	case 3: // start condition: _S2
		goto yystart263
	}

	goto yystate1 // silence unused label error
yystate1:
	c = s.next()
yystart1:
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate4
	case c == '%':
		goto yystate5
	case c == '/':
		goto yystate67
	case c == '0':
		goto yystate72
	case c == '\'':
		goto yystate61
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate3
	case c == '\x00':
		goto yystate2
	case c == '`':
		goto yystate78
	case c >= '1' && c <= '9':
		goto yystate76
	case c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate77
	case c >= 'Â' && c <= 'ß':
		goto yystate79
	case c >= 'à' && c <= 'ï':
		goto yystate81
	case c >= 'ð' && c <= 'ô':
		goto yystate82
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
	goto yyrule81

yystate5:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '%':
		goto yystate6
	case c == 'E' || c == 'e':
		goto yystate7
	case c == 'L' || c == 'l':
		goto yystate20
	case c == 'N' || c == 'n':
		goto yystate24
	case c == 'P' || c == 'p':
		goto yystate32
	case c == 'R' || c == 'r':
		goto yystate36
	case c == 'S' || c == 's':
		goto yystate41
	case c == 'T' || c == 't':
		goto yystate46
	case c == 'U' || c == 'u':
		goto yystate54
	case c == '{':
		goto yystate59
	case c == '}':
		goto yystate60
	}

yystate6:
	c = s.next()
	goto yyrule90

yystate7:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'R' || c == 'r':
		goto yystate8
	}

yystate8:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'R' || c == 'r':
		goto yystate9
	}

yystate9:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'O' || c == 'o':
		goto yystate10
	}

yystate10:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'R' || c == 'r':
		goto yystate11
	}

yystate11:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '-':
		goto yystate12
	}

yystate12:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'V' || c == 'v':
		goto yystate13
	}

yystate13:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'E' || c == 'e':
		goto yystate14
	}

yystate14:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'R' || c == 'r':
		goto yystate15
	}

yystate15:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'B' || c == 'b':
		goto yystate16
	}

yystate16:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'O' || c == 'o':
		goto yystate17
	}

yystate17:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'S' || c == 's':
		goto yystate18
	}

yystate18:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'E' || c == 'e':
		goto yystate19
	}

yystate19:
	c = s.next()
	goto yyrule91

yystate20:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'E' || c == 'e':
		goto yystate21
	}

yystate21:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'F' || c == 'f':
		goto yystate22
	}

yystate22:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'T' || c == 't':
		goto yystate23
	}

yystate23:
	c = s.next()
	goto yyrule92

yystate24:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'O' || c == 'o':
		goto yystate25
	}

yystate25:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'N' || c == 'n':
		goto yystate26
	}

yystate26:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'A' || c == 'a':
		goto yystate27
	}

yystate27:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'S' || c == 's':
		goto yystate28
	}

yystate28:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'S' || c == 's':
		goto yystate29
	}

yystate29:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'O' || c == 'o':
		goto yystate30
	}

yystate30:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'C' || c == 'c':
		goto yystate31
	}

yystate31:
	c = s.next()
	goto yyrule93

yystate32:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'R' || c == 'r':
		goto yystate33
	}

yystate33:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'E' || c == 'e':
		goto yystate34
	}

yystate34:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'C' || c == 'c':
		goto yystate35
	}

yystate35:
	c = s.next()
	goto yyrule94

yystate36:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'I' || c == 'i':
		goto yystate37
	}

yystate37:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'G' || c == 'g':
		goto yystate38
	}

yystate38:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'H' || c == 'h':
		goto yystate39
	}

yystate39:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'T' || c == 't':
		goto yystate40
	}

yystate40:
	c = s.next()
	goto yyrule95

yystate41:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'T' || c == 't':
		goto yystate42
	}

yystate42:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'A' || c == 'a':
		goto yystate43
	}

yystate43:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'R' || c == 'r':
		goto yystate44
	}

yystate44:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'T' || c == 't':
		goto yystate45
	}

yystate45:
	c = s.next()
	goto yyrule96

yystate46:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'O' || c == 'o':
		goto yystate47
	case c == 'Y' || c == 'y':
		goto yystate51
	}

yystate47:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'K' || c == 'k':
		goto yystate48
	}

yystate48:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'E' || c == 'e':
		goto yystate49
	}

yystate49:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'N' || c == 'n':
		goto yystate50
	}

yystate50:
	c = s.next()
	goto yyrule97

yystate51:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'P' || c == 'p':
		goto yystate52
	}

yystate52:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'E' || c == 'e':
		goto yystate53
	}

yystate53:
	c = s.next()
	goto yyrule98

yystate54:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'N' || c == 'n':
		goto yystate55
	}

yystate55:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'I' || c == 'i':
		goto yystate56
	}

yystate56:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'O' || c == 'o':
		goto yystate57
	}

yystate57:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'N' || c == 'n':
		goto yystate58
	}

yystate58:
	c = s.next()
	goto yyrule99

yystate59:
	c = s.next()
	goto yyrule88

yystate60:
	c = s.next()
	goto yyrule89

yystate61:
	c = s.next()
	switch {
	default:
		goto yyrule84
	case c == '\'':
		goto yystate64
	case c == '\\':
		goto yystate65
	case c >= '\x01' && c <= '&' || c >= '(' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate62
	}

yystate62:
	c = s.next()
	switch {
	default:
		goto yyrule84
	case c == '\'':
		goto yystate63
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate62
	}

yystate63:
	c = s.next()
	goto yyrule85

yystate64:
	c = s.next()
	goto yyrule83

yystate65:
	c = s.next()
	switch {
	default:
		goto yyrule84
	case c == '\'':
		goto yystate66
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate62
	}

yystate66:
	c = s.next()
	switch {
	default:
		goto yyrule84
	case c == '\'':
		goto yystate63
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate62
	}

yystate67:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate68
	case c == '/':
		goto yystate71
	}

yystate68:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate69
	case c >= '\x01' && c <= ')' || c >= '+' && c <= 'ÿ':
		goto yystate68
	}

yystate69:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate69
	case c == '/':
		goto yystate70
	case c >= '\x01' && c <= ')' || c >= '+' && c <= '.' || c >= '0' && c <= 'ÿ':
		goto yystate68
	}

yystate70:
	c = s.next()
	goto yyrule3

yystate71:
	c = s.next()
	switch {
	default:
		goto yyrule4
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate71
	}

yystate72:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c == 'X' || c == 'x':
		goto yystate74
	case c >= '0' && c <= '7':
		goto yystate73
	}

yystate73:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c >= '0' && c <= '7':
		goto yystate73
	}

yystate74:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate75
	}

yystate75:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate75
	}

yystate76:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c >= '0' && c <= '9':
		goto yystate76
	}

yystate77:
	c = s.next()
	switch {
	default:
		goto yyrule101
	case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate77
	}

yystate78:
	c = s.next()
	goto yyrule82

yystate79:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate80
	}

yystate80:
	c = s.next()
	goto yyrule103

yystate81:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate79
	}

yystate82:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate81
	}

	goto yystate83 // silence unused label error
yystate83:
	c = s.next()
yystart83:
	switch {
	default:
		goto yyabort
	case c == '!':
		goto yystate84
	case c == '"':
		goto yystate4
	case c == '%':
		goto yystate86
	case c == '&':
		goto yystate88
	case c == '(':
		goto yystate93
	case c == ')':
		goto yystate94
	case c == '*':
		goto yystate95
	case c == '+':
		goto yystate97
	case c == ',':
		goto yystate100
	case c == '-':
		goto yystate101
	case c == '.':
		goto yystate104
	case c == '/':
		goto yystate112
	case c == '0':
		goto yystate114
	case c == ':':
		goto yystate119
	case c == ';':
		goto yystate121
	case c == '<':
		goto yystate122
	case c == '=':
		goto yystate127
	case c == '>':
		goto yystate129
	case c == '[':
		goto yystate134
	case c == '\'':
		goto yystate61
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate3
	case c == '\x00':
		goto yystate2
	case c == ']':
		goto yystate135
	case c == '^':
		goto yystate136
	case c == '`':
		goto yystate78
	case c == 'b':
		goto yystate138
	case c == 'c':
		goto yystate143
	case c == 'd':
		goto yystate159
	case c == 'e':
		goto yystate168
	case c == 'f':
		goto yystate172
	case c == 'g':
		goto yystate188
	case c == 'i':
		goto yystate192
	case c == 'm':
		goto yystate207
	case c == 'p':
		goto yystate210
	case c == 'r':
		goto yystate217
	case c == 's':
		goto yystate227
	case c == 't':
		goto yystate243
	case c == 'v':
		goto yystate247
	case c == '{':
		goto yystate250
	case c == '|':
		goto yystate251
	case c == '}':
		goto yystate254
	case c >= '1' && c <= '9':
		goto yystate118
	case c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'h' || c >= 'j' && c <= 'l' || c == 'n' || c == 'o' || c == 'q' || c == 'u' || c >= 'w' && c <= 'z':
		goto yystate133
	case c >= 'Â' && c <= 'ß':
		goto yystate255
	case c >= 'à' && c <= 'ï':
		goto yystate257
	case c >= 'ð' && c <= 'ô':
		goto yystate258
	}

yystate84:
	c = s.next()
	switch {
	default:
		goto yyrule5
	case c == '=':
		goto yystate85
	}

yystate85:
	c = s.next()
	goto yyrule6

yystate86:
	c = s.next()
	switch {
	default:
		goto yyrule7
	case c == '=':
		goto yystate87
	}

yystate87:
	c = s.next()
	goto yyrule8

yystate88:
	c = s.next()
	switch {
	default:
		goto yyrule9
	case c == '&':
		goto yystate89
	case c == '=':
		goto yystate90
	case c == '^':
		goto yystate91
	}

yystate89:
	c = s.next()
	goto yyrule10

yystate90:
	c = s.next()
	goto yyrule11

yystate91:
	c = s.next()
	switch {
	default:
		goto yyrule12
	case c == '=':
		goto yystate92
	}

yystate92:
	c = s.next()
	goto yyrule13

yystate93:
	c = s.next()
	goto yyrule14

yystate94:
	c = s.next()
	goto yyrule15

yystate95:
	c = s.next()
	switch {
	default:
		goto yyrule16
	case c == '=':
		goto yystate96
	}

yystate96:
	c = s.next()
	goto yyrule17

yystate97:
	c = s.next()
	switch {
	default:
		goto yyrule18
	case c == '+':
		goto yystate98
	case c == '=':
		goto yystate99
	}

yystate98:
	c = s.next()
	goto yyrule19

yystate99:
	c = s.next()
	goto yyrule20

yystate100:
	c = s.next()
	goto yyrule21

yystate101:
	c = s.next()
	switch {
	default:
		goto yyrule22
	case c == '-':
		goto yystate102
	case c == '=':
		goto yystate103
	}

yystate102:
	c = s.next()
	goto yyrule23

yystate103:
	c = s.next()
	goto yyrule24

yystate104:
	c = s.next()
	switch {
	default:
		goto yyrule25
	case c == '.':
		goto yystate105
	case c >= '0' && c <= '9':
		goto yystate107
	}

yystate105:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate106
	}

yystate106:
	c = s.next()
	goto yyrule26

yystate107:
	c = s.next()
	switch {
	default:
		goto yyrule80
	case c == 'E' || c == 'e':
		goto yystate108
	case c == 'i':
		goto yystate111
	case c >= '0' && c <= '9':
		goto yystate107
	}

yystate108:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '+' || c == '-':
		goto yystate109
	case c >= '0' && c <= '9':
		goto yystate110
	}

yystate109:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate110
	}

yystate110:
	c = s.next()
	switch {
	default:
		goto yyrule80
	case c == 'i':
		goto yystate111
	case c >= '0' && c <= '9':
		goto yystate110
	}

yystate111:
	c = s.next()
	goto yyrule78

yystate112:
	c = s.next()
	switch {
	default:
		goto yyrule27
	case c == '*':
		goto yystate68
	case c == '/':
		goto yystate71
	case c == '=':
		goto yystate113
	}

yystate113:
	c = s.next()
	goto yyrule28

yystate114:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c == '.':
		goto yystate107
	case c == '8' || c == '9':
		goto yystate116
	case c == 'E' || c == 'e':
		goto yystate108
	case c == 'X' || c == 'x':
		goto yystate74
	case c == 'i':
		goto yystate117
	case c >= '0' && c <= '7':
		goto yystate115
	}

yystate115:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c == '.':
		goto yystate107
	case c == '8' || c == '9':
		goto yystate116
	case c == 'E' || c == 'e':
		goto yystate108
	case c == 'i':
		goto yystate117
	case c >= '0' && c <= '7':
		goto yystate115
	}

yystate116:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate107
	case c == 'E' || c == 'e':
		goto yystate108
	case c == 'i':
		goto yystate117
	case c >= '0' && c <= '9':
		goto yystate116
	}

yystate117:
	c = s.next()
	goto yyrule77

yystate118:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c == '.':
		goto yystate107
	case c == 'E' || c == 'e':
		goto yystate108
	case c == 'i':
		goto yystate117
	case c >= '0' && c <= '9':
		goto yystate118
	}

yystate119:
	c = s.next()
	switch {
	default:
		goto yyrule29
	case c == '=':
		goto yystate120
	}

yystate120:
	c = s.next()
	goto yyrule30

yystate121:
	c = s.next()
	goto yyrule31

yystate122:
	c = s.next()
	switch {
	default:
		goto yyrule32
	case c == '-':
		goto yystate123
	case c == '<':
		goto yystate124
	case c == '=':
		goto yystate126
	}

yystate123:
	c = s.next()
	goto yyrule33

yystate124:
	c = s.next()
	switch {
	default:
		goto yyrule34
	case c == '=':
		goto yystate125
	}

yystate125:
	c = s.next()
	goto yyrule35

yystate126:
	c = s.next()
	goto yyrule36

yystate127:
	c = s.next()
	switch {
	default:
		goto yyrule37
	case c == '=':
		goto yystate128
	}

yystate128:
	c = s.next()
	goto yyrule38

yystate129:
	c = s.next()
	switch {
	default:
		goto yyrule39
	case c == '=':
		goto yystate130
	case c == '>':
		goto yystate131
	}

yystate130:
	c = s.next()
	goto yyrule40

yystate131:
	c = s.next()
	switch {
	default:
		goto yyrule41
	case c == '=':
		goto yystate132
	}

yystate132:
	c = s.next()
	goto yyrule42

yystate133:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate134:
	c = s.next()
	goto yyrule43

yystate135:
	c = s.next()
	goto yyrule44

yystate136:
	c = s.next()
	switch {
	default:
		goto yyrule45
	case c == '=':
		goto yystate137
	}

yystate137:
	c = s.next()
	goto yyrule46

yystate138:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'r':
		goto yystate139
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate133
	}

yystate139:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate140
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate140:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate141
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate133
	}

yystate141:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'k':
		goto yystate142
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate133
	}

yystate142:
	c = s.next()
	switch {
	default:
		goto yyrule52
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate143:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate144
	case c == 'h':
		goto yystate147
	case c == 'o':
		goto yystate150
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'g' || c >= 'i' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate133
	}

yystate144:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 's':
		goto yystate145
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate133
	}

yystate145:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate146
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate146:
	c = s.next()
	switch {
	default:
		goto yyrule53
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate147:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate148
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate133
	}

yystate148:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'n':
		goto yystate149
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate133
	}

yystate149:
	c = s.next()
	switch {
	default:
		goto yyrule54
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate150:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'n':
		goto yystate151
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate133
	}

yystate151:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 's':
		goto yystate152
	case c == 't':
		goto yystate154
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate152:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 't':
		goto yystate153
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate153:
	c = s.next()
	switch {
	default:
		goto yyrule55
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate154:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'i':
		goto yystate155
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate133
	}

yystate155:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'n':
		goto yystate156
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate133
	}

yystate156:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'u':
		goto yystate157
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate133
	}

yystate157:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate158
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate158:
	c = s.next()
	switch {
	default:
		goto yyrule56
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate159:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate160
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate160:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'f':
		goto yystate161
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate133
	}

yystate161:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate162
	case c == 'e':
		goto yystate166
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate162:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'u':
		goto yystate163
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate133
	}

yystate163:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'l':
		goto yystate164
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate133
	}

yystate164:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 't':
		goto yystate165
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate165:
	c = s.next()
	switch {
	default:
		goto yyrule57
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate166:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'r':
		goto yystate167
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate133
	}

yystate167:
	c = s.next()
	switch {
	default:
		goto yyrule58
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate168:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'l':
		goto yystate169
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate133
	}

yystate169:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 's':
		goto yystate170
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate133
	}

yystate170:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate171
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate171:
	c = s.next()
	switch {
	default:
		goto yyrule59
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate172:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate173
	case c == 'o':
		goto yystate183
	case c == 'u':
		goto yystate185
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'n' || c >= 'p' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate133
	}

yystate173:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'l':
		goto yystate174
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate133
	}

yystate174:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'l':
		goto yystate175
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate133
	}

yystate175:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 't':
		goto yystate176
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate176:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'h':
		goto yystate177
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate133
	}

yystate177:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'r':
		goto yystate178
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate133
	}

yystate178:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'o':
		goto yystate179
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate133
	}

yystate179:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'u':
		goto yystate180
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate133
	}

yystate180:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'g':
		goto yystate181
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate133
	}

yystate181:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'h':
		goto yystate182
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate133
	}

yystate182:
	c = s.next()
	switch {
	default:
		goto yyrule60
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate183:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'r':
		goto yystate184
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate133
	}

yystate184:
	c = s.next()
	switch {
	default:
		goto yyrule61
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate185:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'n':
		goto yystate186
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate133
	}

yystate186:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'c':
		goto yystate187
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate133
	}

yystate187:
	c = s.next()
	switch {
	default:
		goto yyrule62
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate188:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'o':
		goto yystate189
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate133
	}

yystate189:
	c = s.next()
	switch {
	default:
		goto yyrule63
	case c == 't':
		goto yystate190
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate190:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'o':
		goto yystate191
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate133
	}

yystate191:
	c = s.next()
	switch {
	default:
		goto yyrule64
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate192:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'f':
		goto yystate193
	case c == 'm':
		goto yystate194
	case c == 'n':
		goto yystate199
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'l' || c >= 'o' && c <= 'z':
		goto yystate133
	}

yystate193:
	c = s.next()
	switch {
	default:
		goto yyrule65
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate194:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'p':
		goto yystate195
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate133
	}

yystate195:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'o':
		goto yystate196
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate133
	}

yystate196:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'r':
		goto yystate197
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate133
	}

yystate197:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 't':
		goto yystate198
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate198:
	c = s.next()
	switch {
	default:
		goto yyrule66
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate199:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 't':
		goto yystate200
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate200:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate201
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate201:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'r':
		goto yystate202
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate133
	}

yystate202:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'f':
		goto yystate203
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate133
	}

yystate203:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate204
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate133
	}

yystate204:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'c':
		goto yystate205
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate133
	}

yystate205:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate206
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate206:
	c = s.next()
	switch {
	default:
		goto yyrule67
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate207:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate208
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate133
	}

yystate208:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'p':
		goto yystate209
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate133
	}

yystate209:
	c = s.next()
	switch {
	default:
		goto yyrule68
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate210:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate211
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate133
	}

yystate211:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'c':
		goto yystate212
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate133
	}

yystate212:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'k':
		goto yystate213
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate133
	}

yystate213:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate214
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate133
	}

yystate214:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'g':
		goto yystate215
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate133
	}

yystate215:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate216
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate216:
	c = s.next()
	switch {
	default:
		goto yyrule69
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate217:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate218
	case c == 'e':
		goto yystate222
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate218:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'n':
		goto yystate219
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate133
	}

yystate219:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'g':
		goto yystate220
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate133
	}

yystate220:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate221
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate221:
	c = s.next()
	switch {
	default:
		goto yyrule70
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate222:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 't':
		goto yystate223
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate223:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'u':
		goto yystate224
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate133
	}

yystate224:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'r':
		goto yystate225
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate133
	}

yystate225:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'n':
		goto yystate226
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate133
	}

yystate226:
	c = s.next()
	switch {
	default:
		goto yyrule71
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate227:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate228
	case c == 't':
		goto yystate233
	case c == 'w':
		goto yystate238
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 's' || c == 'u' || c == 'v' || c >= 'x' && c <= 'z':
		goto yystate133
	}

yystate228:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'l':
		goto yystate229
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate133
	}

yystate229:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate230
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate230:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'c':
		goto yystate231
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate133
	}

yystate231:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 't':
		goto yystate232
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate232:
	c = s.next()
	switch {
	default:
		goto yyrule72
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate233:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'r':
		goto yystate234
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate133
	}

yystate234:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'u':
		goto yystate235
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate133
	}

yystate235:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'c':
		goto yystate236
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate133
	}

yystate236:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 't':
		goto yystate237
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate237:
	c = s.next()
	switch {
	default:
		goto yyrule73
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate238:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'i':
		goto yystate239
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate133
	}

yystate239:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 't':
		goto yystate240
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate133
	}

yystate240:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'c':
		goto yystate241
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate133
	}

yystate241:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'h':
		goto yystate242
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate133
	}

yystate242:
	c = s.next()
	switch {
	default:
		goto yyrule74
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate243:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'y':
		goto yystate244
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate133
	}

yystate244:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'p':
		goto yystate245
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate133
	}

yystate245:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'e':
		goto yystate246
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate133
	}

yystate246:
	c = s.next()
	switch {
	default:
		goto yyrule75
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate247:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'a':
		goto yystate248
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate133
	}

yystate248:
	c = s.next()
	switch {
	default:
		goto yyrule100
	case c == 'r':
		goto yystate249
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate133
	}

yystate249:
	c = s.next()
	switch {
	default:
		goto yyrule76
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate133
	}

yystate250:
	c = s.next()
	goto yyrule47

yystate251:
	c = s.next()
	switch {
	default:
		goto yyrule48
	case c == '=':
		goto yystate252
	case c == '|':
		goto yystate253
	}

yystate252:
	c = s.next()
	goto yyrule49

yystate253:
	c = s.next()
	goto yyrule50

yystate254:
	c = s.next()
	goto yyrule51

yystate255:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate256
	}

yystate256:
	c = s.next()
	goto yyrule102

yystate257:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate255
	}

yystate258:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate257
	}

	goto yystate259 // silence unused label error
yystate259:
	c = s.next()
yystart259:
	switch {
	default:
		goto yystate260 // c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ'
	case c == '"':
		goto yystate261
	case c == '\\':
		goto yystate262
	case c == '\x00':
		goto yystate2
	}

yystate260:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate261
	case c == '\\':
		goto yystate262
	case c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate260
	}

yystate261:
	c = s.next()
	goto yyrule86

yystate262:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate260
	}

	goto yystate263 // silence unused label error
yystate263:
	c = s.next()
yystart263:
	switch {
	default:
		goto yystate264 // c >= '\x01' && c <= '_' || c >= 'a' && c <= 'ÿ'
	case c == '\x00':
		goto yystate2
	case c == '`':
		goto yystate265
	}

yystate264:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '`':
		goto yystate265
	case c >= '\x01' && c <= '_' || c >= 'a' && c <= 'ÿ':
		goto yystate264
	}

yystate265:
	c = s.next()
	goto yyrule87

yyrule1: // \0
	{
		s.i0++
		return EOF, lval
	}
yyrule2: // [ \t\n\r]+

	goto yystate0
yyrule3: // \/\*([^*]|\*+[^*/])*\*+\/
yyrule4: // \/\/.*
	{
		return COMMENT, string(s.val)
	}
yyrule5: // "!"
	{
		return NOT, lval
	}
yyrule6: // "!="
	{
		return NEQ, lval
	}
yyrule7: // "%"
	{
		return REM, lval
	}
yyrule8: // "%="
	{
		return REM_ASSIGN, lval
	}
yyrule9: // "&"
	{
		return AND, lval
	}
yyrule10: // "&&"
	{
		return LAND, lval
	}
yyrule11: // "&="
	{
		return AND_ASSIGN, lval
	}
yyrule12: // "&^"
	{
		return AND_NOT, lval
	}
yyrule13: // "&^="
	{
		return AND_NOT_ASSIGN, lval
	}
yyrule14: // "("
	{
		return LPAREN, lval
	}
yyrule15: // ")"
	{
		return RPAREN, lval
	}
yyrule16: // "*"
	{
		return MUL, lval
	}
yyrule17: // "*="
	{
		return MUL_ASSIGN, lval
	}
yyrule18: // "+"
	{
		return ADD, lval
	}
yyrule19: // "++"
	{
		return INC, lval
	}
yyrule20: // "+="
	{
		return ADD_ASSIGN, lval
	}
yyrule21: // ","
	{
		return COMMA, lval
	}
yyrule22: // "-"
	{
		return SUB, lval
	}
yyrule23: // "--"
	{
		return DEC, lval
	}
yyrule24: // "-="
	{
		return SUB_ASSIGN, lval
	}
yyrule25: // "."
	{
		return PERIOD, lval
	}
yyrule26: // "..."
	{
		return ELLIPSIS, lval
	}
yyrule27: // "/"
	{
		return QUO, lval
	}
yyrule28: // "/="
	{
		return QUO_ASSIGN, lval
	}
yyrule29: // ":"
	{
		return COLON, lval
	}
yyrule30: // ":="
	{
		return DEFINE, lval
	}
yyrule31: // ";"
	{
		return SEMICOLON, lval
	}
yyrule32: // "<"
	{
		return LSS, lval
	}
yyrule33: // "<-"
	{
		return ARROW, lval
	}
yyrule34: // "<<"
	{
		return SHL, lval
	}
yyrule35: // "<<="
	{
		return SHL_ASSIGN, lval
	}
yyrule36: // "<="
	{
		return LEQ, lval
	}
yyrule37: // "="
	{
		return ASSIGN, lval
	}
yyrule38: // "=="
	{
		return EQL, lval
	}
yyrule39: // ">"
	{
		return GTR, lval
	}
yyrule40: // ">="
	{
		return GEQ, lval
	}
yyrule41: // ">>"
	{
		return SHR, lval
	}
yyrule42: // ">>="
	{
		return SHR_ASSIGN, lval
	}
yyrule43: // "["
	{
		return LBRACK, lval
	}
yyrule44: // "]"
	{
		return RBRACK, lval
	}
yyrule45: // "^"
	{
		return XOR, lval
	}
yyrule46: // "^="
	{
		return XOR_ASSIGN, lval
	}
yyrule47: // "{"
	{
		return LBRACE, lval
	}
yyrule48: // "|"
	{
		return OR, lval
	}
yyrule49: // "|="
	{
		return OR_ASSIGN, lval
	}
yyrule50: // "||"
	{
		return LOR, lval
	}
yyrule51: // "}"
	{
		return RBRACE, lval
	}
yyrule52: // break
	{
		return BREAK, lval
	}
yyrule53: // case
	{
		return CASE, lval
	}
yyrule54: // chan
	{
		return CHAN, lval
	}
yyrule55: // const
	{
		return CONST, lval
	}
yyrule56: // continue
	{
		return CONTINUE, lval
	}
yyrule57: // default
	{
		return DEFAULT, lval
	}
yyrule58: // defer
	{
		return DEFER, lval
	}
yyrule59: // else
	{
		return ELSE, lval
	}
yyrule60: // fallthrough
	{
		return FALLTHROUGH, lval
	}
yyrule61: // for
	{
		return FOR, lval
	}
yyrule62: // func
	{
		return FUNC, lval
	}
yyrule63: // go
	{
		return GO, lval
	}
yyrule64: // goto
	{
		return GOTO, lval
	}
yyrule65: // if
	{
		return IF, lval
	}
yyrule66: // import
	{
		return IMPORT, lval
	}
yyrule67: // interface
	{
		return INTERFACE, lval
	}
yyrule68: // map
	{
		return MAP, lval
	}
yyrule69: // package
	{
		return PACKAGE, lval
	}
yyrule70: // range
	{
		return RANGE, lval
	}
yyrule71: // return
	{
		return RETURN, lval
	}
yyrule72: // select
	{
		return SELECT, lval
	}
yyrule73: // struct
	{
		return STRUCT, lval
	}
yyrule74: // switch
	{
		return SWITCH, lval
	}
yyrule75: // type
	{
		return GO_TYPE, lval
	}
yyrule76: // var
	{
		return VAR, lval
	}
yyrule77: // {imaginary_ilit}
	{
		return s.int(IMAG)
	}
yyrule78: // {imaginary_lit}
	{
		return s.float(IMAG)
	}
yyrule79: // {int_lit}
	{
		return s.int(INT)
	}
yyrule80: // {float_lit}
	{
		return s.float(FLOAT)
	}
yyrule81: // \"
	{
		s.ssc, s.sc = s.sc, _S1
		goto yystate0
	}
yyrule82: // `
	{
		s.ssc, s.sc = s.sc, _S2
		goto yystate0
	}
yyrule83: // ''
yyrule84: // '(\\.)?[^']*
	{
		return CHAR, string(s.val)
	}
yyrule85: // '(\\.)?[^']*'
	{

		s.ssc = s.sc
		if tok, lval = s.str(""); tok != STRING {
			return
		}
		s.i0++
		return CHAR, int32(lval.(string)[0])
	}
yyrule86: // (\\.|[^\\"])*\"
	{
		return s.str("\"")
	}
yyrule87: // ([^`]|\n)*`
	{
		return s.str("`")
	}
yyrule88: // "%{"
	{
		return LCURL, lval
	}
yyrule89: // "%}"
	{
		return RCURL, lval
	}
yyrule90: // %%
	{
		return MARK, lval
	}
yyrule91: // %{errorverbose}
	{
		return ERR_VERBOSE, lval
	}
yyrule92: // %{left}
	{
		return LEFT, lval
	}
yyrule93: // %{nonassoc}
	{
		return NONASSOC, lval
	}
yyrule94: // %{prec}
	{
		return PREC, lval
	}
yyrule95: // %{right}
	{
		return RIGHT, lval
	}
yyrule96: // %{start}
	{
		return START, lval
	}
yyrule97: // %{token}
	{
		return TOKEN, lval
	}
yyrule98: // %{type}
	{
		return TYPE, lval
	}
yyrule99: // %{union}
	{
		return UNION, lval
	}
yyrule100: // [a-zA-Z_][a-zA-Z0-9_]*
	{

		if c >= '\xC2' && c <= '\xF4' {
			s.i--
			s.NCol--
			for rune := rune(1); rune >= 0; rune = s.getRune(true, false) {
				tok, lval = IDENT, string(s.src[s.i0-1:s.i])
			}
			s.next()
			return
		}
		return IDENT, string(s.val)
	}
yyrule101: // [a-zA-Z_][a-zA-Z0-9_.]*
	{

		if c >= '\xC2' && c <= '\xF4' {
			s.i--
			s.NCol--
			for rune := rune(1); rune >= 0; rune = s.getRune(true, true) {
				tok, lval = IDENTIFIER, string(s.src[s.i0-1:s.i])
			}
			s.next()
			return
		}
		return IDENTIFIER, string(s.val)
	}
yyrule102: // {non_ascii}
	{

		s.i = s.i0 - 1
		if rune := s.getRune(false, false); rune < 0 {
			_, sz := utf8.DecodeRune(s.src[s.i:])
			s.i += sz
			s.next()
			s.err("expected unicode lettter, got %U", rune)
			return ILLEGAL, string(-rune)
		}
		for rune := rune(1); rune >= 0; rune = s.getRune(true, false) {
		}
		s.next()
		return IDENT, string(s.src[s.i0-1 : s.i-1])
	}
yyrule103: // {non_ascii}
	{

		s.i = s.i0 - 1
		if rune := s.getRune(false, false); rune < 0 {
			_, sz := utf8.DecodeRune(s.src[s.i:])
			s.i += sz
			s.next()
			s.err("expected unicode lettter, got %U", rune)
			return ILLEGAL, string(-rune)
		}
		for rune := rune(1); rune >= 0; rune = s.getRune(true, true) {
		}
		s.next()
		return IDENTIFIER, string(s.src[s.i0-1 : s.i-1])
	}
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	s.next()
	return ILLEGAL, string(c0)
}

func (s *Scanner) getRune(acceptDigits, acceptDot bool) (r rune) {
	var sz int
	if r, sz = utf8.DecodeRune(s.src[s.i:]); sz != 0 &&
		(r == '_' || unicode.IsLetter(r) || (acceptDigits && unicode.IsDigit(r)) || (acceptDot && r == '.')) {
		s.i += sz
		s.NCol += sz
		return
	}

	return -r
}

func (s *Scanner) str(pref string) (tok Token, lval interface{}) {
	s.sc = s.ssc
	ss := pref + string(s.val)
	ss, err := strconv.Unquote(ss)
	if err != nil {
		s.err("string literal %q: %v", ss, err)
		return ILLEGAL, ss
	}

	s.i0--
	return STRING, ss
}

func (s *Scanner) int(tk Token) (tok Token, lval interface{}) {
	tok = tk
	if tok == IMAG {
		s.val = s.val[:len(s.val)-1]
	}
	n, err := strconv.ParseUint(string(s.val), 0, 64)
	if err != nil {
		lval = string(s.val)
	} else if tok == IMAG {
		lval = complex(0, float64(n))
	} else {
		lval = n
	}
	return
}

func (s *Scanner) float(tk Token) (tok Token, lval interface{}) {
	tok = tk
	if tok == IMAG {
		s.val = s.val[:len(s.val)-1]
	}
	n, err := strconv.ParseFloat(string(s.val), 64)
	if err != nil {
		lval = string(s.val)
	} else if tok == IMAG {
		lval = complex(0, n)
	} else {
		lval = n
	}
	return
}
