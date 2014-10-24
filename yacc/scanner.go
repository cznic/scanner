// Copyright (c) 2014 The scanner Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// CAUTION: If this file is a Go source file (*.go), it was generated
// automatically by '$ golex' from a *.l file - DO NOT EDIT in that case!

// Package scanner implements a scanner for yacc[1] source text with actions
// written in Go. It takes a []byte as source which can then be tokenized
// through repeated calls to the Scan method.
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
	"go/token"
	"strconv"
	"strings"
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
	_ Token = 0xE000 + iota

	// ------------------------------------------- Go mode or shared tokens

	// Special tokens
	ILLEGAL
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

	DLR_DLR     // $$
	DLR_NUM     // $1
	DLR_TAG_DLR // $<tag>$
	DLR_TAG_NUM // $<tag>2

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

	DLR_DLR:     "DLR_DLR",
	DLR_NUM:     "DLR_NUM",
	DLR_TAG_DLR: "DLR_TAG_DLR",
	DLR_TAG_NUM: "DLR_TAG_NUM",

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

	return fmt.Sprintf("Token(%d)", int(i))
}

// A Scanner holds the scanner's internal state while processing a given text.
type Scanner struct {
	Col    int     // Starting column of the last scanned token.
	Errors []error // List of accumulated errors.
	FName  string  // File name (reported) of the scanned source.
	Line   int     // Starting line of the last scanned token.
	NCol   int     // Starting column (reported) for the next scanned token.
	NLine  int     // Starting line (reported) for the next scanned token.
	c      int
	file   *token.File
	i      int
	i0     int
	sc     int
	src    []byte
	ssc    int // saved state condition
	val    []byte
}

// New returns a newly created Scanner and set its FName to fname
func New(fset *token.FileSet, fname string, src []byte) (s *Scanner) {
	if len(src) > 2 && src[0] == 0xEF && src[1] == 0xBB && src[2] == 0xBF {
		src = src[3:]
	}
	s = &Scanner{
		FName: fname,
		src:   src,
		NLine: 1,
		NCol:  0,
	}
	s.file = fset.AddFile(fname, -1, len(src))
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
		s.file.AddLine(s.i)
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
	err := fmt.Errorf(fmt.Sprintf("%s:%d:%d ", s.FName, s.Line, s.Col)+format, arg...)
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
func (s *Scanner) Scan() (tok Token, lval interface{}, num int) {
	tok, lval, num = s.ScanRaw()
	if s.sc != _INITIAL || tok != IDENTIFIER {
		return
	}

	idLine, idCol := s.Line, s.Col
	i, nl, nc, c := s.i, s.NLine, s.NCol, s.c
	i0 := s.i0
	tok2, lit, _ := s.ScanRaw()
	s.i0 = i0
	s.Line, s.Col = idLine, idCol
	if tok2 == ILLEGAL && lit.(string) == ":" {
		return C_IDENTIFIER, lval, 0
	}

	s.i, s.NLine, s.NCol, s.c = i, nl, nc, c
	return
}

// ScanRaw scans the next token and returns the token and its value if
// applicable.  The source end is indicated by EOF.
//
// If the returned token is IDENT, INT, FLOAT, IMAG, CHAR, STRING or COMMENT,
// lval has has the corresponding value - not the string representation of the
// value. However, numeric literals which overflow the corresponding Go
// predeclared types are returned as string.
//
// If the returned token is DLR_NUM or DLR_TAG_DLR, num has the corresponding
// number (int) or lval has the corresponding tag (string).
//
// If the returned token is DLR_TAG_NUM, lval has the corresponding tag (string)
// and num has the corresponding number.
//
// If the returned token is ILLEGAL, the literal string is the offending
// character or number/string/char literal.
func (s *Scanner) ScanRaw() (tok Token, lval interface{}, num int) {
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
		goto yystart84
	case 2: // start condition: _S1
		goto yystart270
	case 3: // start condition: _S2
		goto yystart274
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
	case c == ',':
		goto yystate67
	case c == '/':
		goto yystate68
	case c == '0':
		goto yystate73
	case c == '\'':
		goto yystate61
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate3
	case c == '\x00':
		goto yystate2
	case c == '`':
		goto yystate79
	case c >= '1' && c <= '9':
		goto yystate77
	case c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate78
	case c >= 'Â' && c <= 'ß':
		goto yystate80
	case c >= 'à' && c <= 'ï':
		goto yystate82
	case c >= 'ð' && c <= 'ô':
		goto yystate83
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
	goto yyrule85

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
	goto yyrule94

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
	goto yyrule95

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
	goto yyrule96

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
	goto yyrule97

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
	goto yyrule98

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
	goto yyrule99

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
	goto yyrule100

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
	goto yyrule101

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
	goto yyrule102

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
	goto yyrule103

yystate59:
	c = s.next()
	goto yyrule92

yystate60:
	c = s.next()
	goto yyrule93

yystate61:
	c = s.next()
	switch {
	default:
		goto yyrule88
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
		goto yyrule88
	case c == '\'':
		goto yystate63
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate62
	}

yystate63:
	c = s.next()
	goto yyrule89

yystate64:
	c = s.next()
	goto yyrule87

yystate65:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == '\'':
		goto yystate66
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate62
	}

yystate66:
	c = s.next()
	switch {
	default:
		goto yyrule88
	case c == '\'':
		goto yystate63
	case c >= '\x01' && c <= '&' || c >= '(' && c <= 'ÿ':
		goto yystate62
	}

yystate67:
	c = s.next()
	goto yyrule104

yystate68:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate69
	case c == '/':
		goto yystate72
	}

yystate69:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate70
	case c >= '\x01' && c <= ')' || c >= '+' && c <= 'ÿ':
		goto yystate69
	}

yystate70:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '*':
		goto yystate70
	case c == '/':
		goto yystate71
	case c >= '\x01' && c <= ')' || c >= '+' && c <= '.' || c >= '0' && c <= 'ÿ':
		goto yystate69
	}

yystate71:
	c = s.next()
	goto yyrule3

yystate72:
	c = s.next()
	switch {
	default:
		goto yyrule4
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate72
	}

yystate73:
	c = s.next()
	switch {
	default:
		goto yyrule83
	case c == 'X' || c == 'x':
		goto yystate75
	case c >= '0' && c <= '7':
		goto yystate74
	}

yystate74:
	c = s.next()
	switch {
	default:
		goto yyrule83
	case c >= '0' && c <= '7':
		goto yystate74
	}

yystate75:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate76
	}

yystate76:
	c = s.next()
	switch {
	default:
		goto yyrule83
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate76
	}

yystate77:
	c = s.next()
	switch {
	default:
		goto yyrule83
	case c >= '0' && c <= '9':
		goto yystate77
	}

yystate78:
	c = s.next()
	switch {
	default:
		goto yyrule106
	case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate78
	}

yystate79:
	c = s.next()
	goto yyrule86

yystate80:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate81
	}

yystate81:
	c = s.next()
	goto yyrule108

yystate82:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate80
	}

yystate83:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate82
	}

	goto yystate84 // silence unused label error
yystate84:
	c = s.next()
yystart84:
	switch {
	default:
		goto yyabort
	case c == '!':
		goto yystate85
	case c == '"':
		goto yystate4
	case c == '$':
		goto yystate87
	case c == '%':
		goto yystate97
	case c == '&':
		goto yystate99
	case c == '(':
		goto yystate104
	case c == ')':
		goto yystate105
	case c == '*':
		goto yystate106
	case c == '+':
		goto yystate108
	case c == ',':
		goto yystate111
	case c == '-':
		goto yystate112
	case c == '.':
		goto yystate115
	case c == '/':
		goto yystate123
	case c == '0':
		goto yystate125
	case c == ':':
		goto yystate130
	case c == ';':
		goto yystate132
	case c == '<':
		goto yystate133
	case c == '=':
		goto yystate138
	case c == '>':
		goto yystate140
	case c == '[':
		goto yystate145
	case c == '\'':
		goto yystate61
	case c == '\t' || c == '\n' || c == '\r' || c == ' ':
		goto yystate3
	case c == '\x00':
		goto yystate2
	case c == ']':
		goto yystate146
	case c == '^':
		goto yystate147
	case c == '`':
		goto yystate79
	case c == 'b':
		goto yystate149
	case c == 'c':
		goto yystate154
	case c == 'd':
		goto yystate170
	case c == 'e':
		goto yystate179
	case c == 'f':
		goto yystate183
	case c == 'g':
		goto yystate199
	case c == 'i':
		goto yystate203
	case c == 'm':
		goto yystate218
	case c == 'p':
		goto yystate221
	case c == 'r':
		goto yystate228
	case c == 's':
		goto yystate238
	case c == 't':
		goto yystate254
	case c == 'v':
		goto yystate258
	case c == '{':
		goto yystate261
	case c == '|':
		goto yystate262
	case c == '}':
		goto yystate265
	case c >= '1' && c <= '9':
		goto yystate129
	case c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'h' || c >= 'j' && c <= 'l' || c == 'n' || c == 'o' || c == 'q' || c == 'u' || c >= 'w' && c <= 'z':
		goto yystate144
	case c >= 'Â' && c <= 'ß':
		goto yystate266
	case c >= 'à' && c <= 'ï':
		goto yystate268
	case c >= 'ð' && c <= 'ô':
		goto yystate269
	}

yystate85:
	c = s.next()
	switch {
	default:
		goto yyrule9
	case c == '=':
		goto yystate86
	}

yystate86:
	c = s.next()
	goto yyrule10

yystate87:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '$':
		goto yystate88
	case c == '-':
		goto yystate89
	case c == '<':
		goto yystate91
	case c >= '0' && c <= '9':
		goto yystate90
	}

yystate88:
	c = s.next()
	goto yyrule8

yystate89:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate90
	}

yystate90:
	c = s.next()
	switch {
	default:
		goto yyrule5
	case c >= '0' && c <= '9':
		goto yystate90
	}

yystate91:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate92
	}

yystate92:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '.' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate92
	case c == '>':
		goto yystate93
	}

yystate93:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '$':
		goto yystate94
	case c == '-':
		goto yystate95
	case c >= '0' && c <= '9':
		goto yystate96
	}

yystate94:
	c = s.next()
	goto yyrule6

yystate95:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate96
	}

yystate96:
	c = s.next()
	switch {
	default:
		goto yyrule7
	case c >= '0' && c <= '9':
		goto yystate96
	}

yystate97:
	c = s.next()
	switch {
	default:
		goto yyrule11
	case c == '=':
		goto yystate98
	}

yystate98:
	c = s.next()
	goto yyrule12

yystate99:
	c = s.next()
	switch {
	default:
		goto yyrule13
	case c == '&':
		goto yystate100
	case c == '=':
		goto yystate101
	case c == '^':
		goto yystate102
	}

yystate100:
	c = s.next()
	goto yyrule14

yystate101:
	c = s.next()
	goto yyrule15

yystate102:
	c = s.next()
	switch {
	default:
		goto yyrule16
	case c == '=':
		goto yystate103
	}

yystate103:
	c = s.next()
	goto yyrule17

yystate104:
	c = s.next()
	goto yyrule18

yystate105:
	c = s.next()
	goto yyrule19

yystate106:
	c = s.next()
	switch {
	default:
		goto yyrule20
	case c == '=':
		goto yystate107
	}

yystate107:
	c = s.next()
	goto yyrule21

yystate108:
	c = s.next()
	switch {
	default:
		goto yyrule22
	case c == '+':
		goto yystate109
	case c == '=':
		goto yystate110
	}

yystate109:
	c = s.next()
	goto yyrule23

yystate110:
	c = s.next()
	goto yyrule24

yystate111:
	c = s.next()
	goto yyrule25

yystate112:
	c = s.next()
	switch {
	default:
		goto yyrule26
	case c == '-':
		goto yystate113
	case c == '=':
		goto yystate114
	}

yystate113:
	c = s.next()
	goto yyrule27

yystate114:
	c = s.next()
	goto yyrule28

yystate115:
	c = s.next()
	switch {
	default:
		goto yyrule29
	case c == '.':
		goto yystate116
	case c >= '0' && c <= '9':
		goto yystate118
	}

yystate116:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate117
	}

yystate117:
	c = s.next()
	goto yyrule30

yystate118:
	c = s.next()
	switch {
	default:
		goto yyrule84
	case c == 'E' || c == 'e':
		goto yystate119
	case c == 'i':
		goto yystate122
	case c >= '0' && c <= '9':
		goto yystate118
	}

yystate119:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '+' || c == '-':
		goto yystate120
	case c >= '0' && c <= '9':
		goto yystate121
	}

yystate120:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9':
		goto yystate121
	}

yystate121:
	c = s.next()
	switch {
	default:
		goto yyrule84
	case c == 'i':
		goto yystate122
	case c >= '0' && c <= '9':
		goto yystate121
	}

yystate122:
	c = s.next()
	goto yyrule82

yystate123:
	c = s.next()
	switch {
	default:
		goto yyrule31
	case c == '*':
		goto yystate69
	case c == '/':
		goto yystate72
	case c == '=':
		goto yystate124
	}

yystate124:
	c = s.next()
	goto yyrule32

yystate125:
	c = s.next()
	switch {
	default:
		goto yyrule83
	case c == '.':
		goto yystate118
	case c == '8' || c == '9':
		goto yystate127
	case c == 'E' || c == 'e':
		goto yystate119
	case c == 'X' || c == 'x':
		goto yystate75
	case c == 'i':
		goto yystate128
	case c >= '0' && c <= '7':
		goto yystate126
	}

yystate126:
	c = s.next()
	switch {
	default:
		goto yyrule83
	case c == '.':
		goto yystate118
	case c == '8' || c == '9':
		goto yystate127
	case c == 'E' || c == 'e':
		goto yystate119
	case c == 'i':
		goto yystate128
	case c >= '0' && c <= '7':
		goto yystate126
	}

yystate127:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '.':
		goto yystate118
	case c == 'E' || c == 'e':
		goto yystate119
	case c == 'i':
		goto yystate128
	case c >= '0' && c <= '9':
		goto yystate127
	}

yystate128:
	c = s.next()
	goto yyrule81

yystate129:
	c = s.next()
	switch {
	default:
		goto yyrule83
	case c == '.':
		goto yystate118
	case c == 'E' || c == 'e':
		goto yystate119
	case c == 'i':
		goto yystate128
	case c >= '0' && c <= '9':
		goto yystate129
	}

yystate130:
	c = s.next()
	switch {
	default:
		goto yyrule33
	case c == '=':
		goto yystate131
	}

yystate131:
	c = s.next()
	goto yyrule34

yystate132:
	c = s.next()
	goto yyrule35

yystate133:
	c = s.next()
	switch {
	default:
		goto yyrule36
	case c == '-':
		goto yystate134
	case c == '<':
		goto yystate135
	case c == '=':
		goto yystate137
	}

yystate134:
	c = s.next()
	goto yyrule37

yystate135:
	c = s.next()
	switch {
	default:
		goto yyrule38
	case c == '=':
		goto yystate136
	}

yystate136:
	c = s.next()
	goto yyrule39

yystate137:
	c = s.next()
	goto yyrule40

yystate138:
	c = s.next()
	switch {
	default:
		goto yyrule41
	case c == '=':
		goto yystate139
	}

yystate139:
	c = s.next()
	goto yyrule42

yystate140:
	c = s.next()
	switch {
	default:
		goto yyrule43
	case c == '=':
		goto yystate141
	case c == '>':
		goto yystate142
	}

yystate141:
	c = s.next()
	goto yyrule44

yystate142:
	c = s.next()
	switch {
	default:
		goto yyrule45
	case c == '=':
		goto yystate143
	}

yystate143:
	c = s.next()
	goto yyrule46

yystate144:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate145:
	c = s.next()
	goto yyrule47

yystate146:
	c = s.next()
	goto yyrule48

yystate147:
	c = s.next()
	switch {
	default:
		goto yyrule49
	case c == '=':
		goto yystate148
	}

yystate148:
	c = s.next()
	goto yyrule50

yystate149:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'r':
		goto yystate150
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate144
	}

yystate150:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate151
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate151:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate152
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate144
	}

yystate152:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'k':
		goto yystate153
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate144
	}

yystate153:
	c = s.next()
	switch {
	default:
		goto yyrule56
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate154:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate155
	case c == 'h':
		goto yystate158
	case c == 'o':
		goto yystate161
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'g' || c >= 'i' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate144
	}

yystate155:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 's':
		goto yystate156
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate144
	}

yystate156:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate157
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate157:
	c = s.next()
	switch {
	default:
		goto yyrule57
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate158:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate159
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate144
	}

yystate159:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'n':
		goto yystate160
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate144
	}

yystate160:
	c = s.next()
	switch {
	default:
		goto yyrule58
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate161:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'n':
		goto yystate162
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate144
	}

yystate162:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 's':
		goto yystate163
	case c == 't':
		goto yystate165
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate163:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 't':
		goto yystate164
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate164:
	c = s.next()
	switch {
	default:
		goto yyrule59
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate165:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'i':
		goto yystate166
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate144
	}

yystate166:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'n':
		goto yystate167
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate144
	}

yystate167:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'u':
		goto yystate168
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate144
	}

yystate168:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate169
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate169:
	c = s.next()
	switch {
	default:
		goto yyrule60
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate170:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate171
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate171:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'f':
		goto yystate172
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate144
	}

yystate172:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate173
	case c == 'e':
		goto yystate177
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate173:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'u':
		goto yystate174
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate144
	}

yystate174:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'l':
		goto yystate175
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate144
	}

yystate175:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 't':
		goto yystate176
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate176:
	c = s.next()
	switch {
	default:
		goto yyrule61
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate177:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'r':
		goto yystate178
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate144
	}

yystate178:
	c = s.next()
	switch {
	default:
		goto yyrule62
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate179:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'l':
		goto yystate180
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate144
	}

yystate180:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 's':
		goto yystate181
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'r' || c >= 't' && c <= 'z':
		goto yystate144
	}

yystate181:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate182
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate182:
	c = s.next()
	switch {
	default:
		goto yyrule63
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate183:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate184
	case c == 'o':
		goto yystate194
	case c == 'u':
		goto yystate196
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'n' || c >= 'p' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate144
	}

yystate184:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'l':
		goto yystate185
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate144
	}

yystate185:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'l':
		goto yystate186
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate144
	}

yystate186:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 't':
		goto yystate187
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate187:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'h':
		goto yystate188
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate144
	}

yystate188:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'r':
		goto yystate189
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate144
	}

yystate189:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'o':
		goto yystate190
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate144
	}

yystate190:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'u':
		goto yystate191
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate144
	}

yystate191:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'g':
		goto yystate192
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate144
	}

yystate192:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'h':
		goto yystate193
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate144
	}

yystate193:
	c = s.next()
	switch {
	default:
		goto yyrule64
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate194:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'r':
		goto yystate195
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate144
	}

yystate195:
	c = s.next()
	switch {
	default:
		goto yyrule65
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate196:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'n':
		goto yystate197
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate144
	}

yystate197:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'c':
		goto yystate198
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate144
	}

yystate198:
	c = s.next()
	switch {
	default:
		goto yyrule66
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate199:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'o':
		goto yystate200
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate144
	}

yystate200:
	c = s.next()
	switch {
	default:
		goto yyrule67
	case c == 't':
		goto yystate201
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate201:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'o':
		goto yystate202
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate144
	}

yystate202:
	c = s.next()
	switch {
	default:
		goto yyrule68
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate203:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'f':
		goto yystate204
	case c == 'm':
		goto yystate205
	case c == 'n':
		goto yystate210
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'l' || c >= 'o' && c <= 'z':
		goto yystate144
	}

yystate204:
	c = s.next()
	switch {
	default:
		goto yyrule69
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate205:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'p':
		goto yystate206
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate144
	}

yystate206:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'o':
		goto yystate207
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'n' || c >= 'p' && c <= 'z':
		goto yystate144
	}

yystate207:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'r':
		goto yystate208
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate144
	}

yystate208:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 't':
		goto yystate209
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate209:
	c = s.next()
	switch {
	default:
		goto yyrule70
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate210:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 't':
		goto yystate211
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate211:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate212
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate212:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'r':
		goto yystate213
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate144
	}

yystate213:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'f':
		goto yystate214
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'e' || c >= 'g' && c <= 'z':
		goto yystate144
	}

yystate214:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate215
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate144
	}

yystate215:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'c':
		goto yystate216
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate144
	}

yystate216:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate217
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate217:
	c = s.next()
	switch {
	default:
		goto yyrule71
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate218:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate219
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate144
	}

yystate219:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'p':
		goto yystate220
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate144
	}

yystate220:
	c = s.next()
	switch {
	default:
		goto yyrule72
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate221:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate222
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate144
	}

yystate222:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'c':
		goto yystate223
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate144
	}

yystate223:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'k':
		goto yystate224
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'j' || c >= 'l' && c <= 'z':
		goto yystate144
	}

yystate224:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate225
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate144
	}

yystate225:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'g':
		goto yystate226
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate144
	}

yystate226:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate227
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate227:
	c = s.next()
	switch {
	default:
		goto yyrule73
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate228:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate229
	case c == 'e':
		goto yystate233
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate229:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'n':
		goto yystate230
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate144
	}

yystate230:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'g':
		goto yystate231
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'f' || c >= 'h' && c <= 'z':
		goto yystate144
	}

yystate231:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate232
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate232:
	c = s.next()
	switch {
	default:
		goto yyrule74
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate233:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 't':
		goto yystate234
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate234:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'u':
		goto yystate235
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate144
	}

yystate235:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'r':
		goto yystate236
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate144
	}

yystate236:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'n':
		goto yystate237
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'm' || c >= 'o' && c <= 'z':
		goto yystate144
	}

yystate237:
	c = s.next()
	switch {
	default:
		goto yyrule75
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate238:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate239
	case c == 't':
		goto yystate244
	case c == 'w':
		goto yystate249
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 's' || c == 'u' || c == 'v' || c >= 'x' && c <= 'z':
		goto yystate144
	}

yystate239:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'l':
		goto yystate240
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'k' || c >= 'm' && c <= 'z':
		goto yystate144
	}

yystate240:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate241
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate241:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'c':
		goto yystate242
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate144
	}

yystate242:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 't':
		goto yystate243
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate243:
	c = s.next()
	switch {
	default:
		goto yyrule76
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate244:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'r':
		goto yystate245
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate144
	}

yystate245:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'u':
		goto yystate246
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 't' || c >= 'v' && c <= 'z':
		goto yystate144
	}

yystate246:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'c':
		goto yystate247
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate144
	}

yystate247:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 't':
		goto yystate248
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate248:
	c = s.next()
	switch {
	default:
		goto yyrule77
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate249:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'i':
		goto yystate250
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'h' || c >= 'j' && c <= 'z':
		goto yystate144
	}

yystate250:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 't':
		goto yystate251
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 's' || c >= 'u' && c <= 'z':
		goto yystate144
	}

yystate251:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'c':
		goto yystate252
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c == 'a' || c == 'b' || c >= 'd' && c <= 'z':
		goto yystate144
	}

yystate252:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'h':
		goto yystate253
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'g' || c >= 'i' && c <= 'z':
		goto yystate144
	}

yystate253:
	c = s.next()
	switch {
	default:
		goto yyrule78
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate254:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'y':
		goto yystate255
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'x' || c == 'z':
		goto yystate144
	}

yystate255:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'p':
		goto yystate256
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'o' || c >= 'q' && c <= 'z':
		goto yystate144
	}

yystate256:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'e':
		goto yystate257
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'd' || c >= 'f' && c <= 'z':
		goto yystate144
	}

yystate257:
	c = s.next()
	switch {
	default:
		goto yyrule79
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate258:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'a':
		goto yystate259
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'b' && c <= 'z':
		goto yystate144
	}

yystate259:
	c = s.next()
	switch {
	default:
		goto yyrule105
	case c == 'r':
		goto yystate260
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'q' || c >= 's' && c <= 'z':
		goto yystate144
	}

yystate260:
	c = s.next()
	switch {
	default:
		goto yyrule80
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate144
	}

yystate261:
	c = s.next()
	goto yyrule51

yystate262:
	c = s.next()
	switch {
	default:
		goto yyrule52
	case c == '=':
		goto yystate263
	case c == '|':
		goto yystate264
	}

yystate263:
	c = s.next()
	goto yyrule53

yystate264:
	c = s.next()
	goto yyrule54

yystate265:
	c = s.next()
	goto yyrule55

yystate266:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate267
	}

yystate267:
	c = s.next()
	goto yyrule107

yystate268:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate266
	}

yystate269:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate268
	}

	goto yystate270 // silence unused label error
yystate270:
	c = s.next()
yystart270:
	switch {
	default:
		goto yystate271 // c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ'
	case c == '"':
		goto yystate272
	case c == '\\':
		goto yystate273
	case c == '\x00':
		goto yystate2
	}

yystate271:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate272
	case c == '\\':
		goto yystate273
	case c >= '\x01' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate271
	}

yystate272:
	c = s.next()
	goto yyrule90

yystate273:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate271
	}

	goto yystate274 // silence unused label error
yystate274:
	c = s.next()
yystart274:
	switch {
	default:
		goto yystate275 // c >= '\x01' && c <= '_' || c >= 'a' && c <= 'ÿ'
	case c == '\x00':
		goto yystate2
	case c == '`':
		goto yystate276
	}

yystate275:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '`':
		goto yystate276
	case c >= '\x01' && c <= '_' || c >= 'a' && c <= 'ÿ':
		goto yystate275
	}

yystate276:
	c = s.next()
	goto yyrule91

yyrule1: // \0
	{
		s.i0++
		return EOF, lval, 0
	}
yyrule2: // [ \t\n\r]+

	goto yystate0
yyrule3: // \/\*([^*]|\*+[^*/])*\*+\/
yyrule4: // \/\/.*
	{
		return COMMENT, string(s.val), 0
	}
yyrule5: // $-?{D}
	{

		v := s.val[1:]
		neg := v[0] == '-'
		if neg {
			v = v[1:]
		}
		n, err := strconv.ParseUint(string(v), 0, 32)
		if err != nil {
			fmt.Printf("%q\n", v)
			return ILLEGAL, nil, 0
		}
		num = int(n)
		if neg {
			num *= -1
		}
		return DLR_NUM, lval, num
	}
yyrule6: // $<[a-zA-Z_][a-zA-Z0-9_.]*>\$
	{

		return DLR_TAG_DLR, string(s.val[2 : len(s.val)-2]), 0
	}
yyrule7: // $<[a-zA-Z_][a-zA-Z0-9_.]*>-?{D}
	{

		a := strings.Split(string(s.val[2:]), ">")
		v := a[1]
		neg := v[0] == '-'
		if neg {
			v = v[1:]
		}
		n, err := strconv.ParseUint(string(v), 0, 32)
		if err != nil {
			fmt.Printf("%q\n", v)
			return ILLEGAL, nil, 0
		}
		num = int(n)
		if neg {
			num *= -1
		}
		return DLR_TAG_NUM, a[0], num
	}
yyrule8: // "$$"
	{
		return DLR_DLR, lval, 0
	}
yyrule9: // "!"
	{
		return NOT, lval, 0
	}
yyrule10: // "!="
	{
		return NEQ, lval, 0
	}
yyrule11: // "%"
	{
		return REM, lval, 0
	}
yyrule12: // "%="
	{
		return REM_ASSIGN, lval, 0
	}
yyrule13: // "&"
	{
		return AND, lval, 0
	}
yyrule14: // "&&"
	{
		return LAND, lval, 0
	}
yyrule15: // "&="
	{
		return AND_ASSIGN, lval, 0
	}
yyrule16: // "&^"
	{
		return AND_NOT, lval, 0
	}
yyrule17: // "&^="
	{
		return AND_NOT_ASSIGN, lval, 0
	}
yyrule18: // "("
	{
		return LPAREN, lval, 0
	}
yyrule19: // ")"
	{
		return RPAREN, lval, 0
	}
yyrule20: // "*"
	{
		return MUL, lval, 0
	}
yyrule21: // "*="
	{
		return MUL_ASSIGN, lval, 0
	}
yyrule22: // "+"
	{
		return ADD, lval, 0
	}
yyrule23: // "++"
	{
		return INC, lval, 0
	}
yyrule24: // "+="
	{
		return ADD_ASSIGN, lval, 0
	}
yyrule25: // ","
	{
		return COMMA, lval, 0
	}
yyrule26: // "-"
	{
		return SUB, lval, 0
	}
yyrule27: // "--"
	{
		return DEC, lval, 0
	}
yyrule28: // "-="
	{
		return SUB_ASSIGN, lval, 0
	}
yyrule29: // "."
	{
		return PERIOD, lval, 0
	}
yyrule30: // "..."
	{
		return ELLIPSIS, lval, 0
	}
yyrule31: // "/"
	{
		return QUO, lval, 0
	}
yyrule32: // "/="
	{
		return QUO_ASSIGN, lval, 0
	}
yyrule33: // ":"
	{
		return COLON, lval, 0
	}
yyrule34: // ":="
	{
		return DEFINE, lval, 0
	}
yyrule35: // ";"
	{
		return SEMICOLON, lval, 0
	}
yyrule36: // "<"
	{
		return LSS, lval, 0
	}
yyrule37: // "<-"
	{
		return ARROW, lval, 0
	}
yyrule38: // "<<"
	{
		return SHL, lval, 0
	}
yyrule39: // "<<="
	{
		return SHL_ASSIGN, lval, 0
	}
yyrule40: // "<="
	{
		return LEQ, lval, 0
	}
yyrule41: // "="
	{
		return ASSIGN, lval, 0
	}
yyrule42: // "=="
	{
		return EQL, lval, 0
	}
yyrule43: // ">"
	{
		return GTR, lval, 0
	}
yyrule44: // ">="
	{
		return GEQ, lval, 0
	}
yyrule45: // ">>"
	{
		return SHR, lval, 0
	}
yyrule46: // ">>="
	{
		return SHR_ASSIGN, lval, 0
	}
yyrule47: // "["
	{
		return LBRACK, lval, 0
	}
yyrule48: // "]"
	{
		return RBRACK, lval, 0
	}
yyrule49: // "^"
	{
		return XOR, lval, 0
	}
yyrule50: // "^="
	{
		return XOR_ASSIGN, lval, 0
	}
yyrule51: // "{"
	{
		return LBRACE, lval, 0
	}
yyrule52: // "|"
	{
		return OR, lval, 0
	}
yyrule53: // "|="
	{
		return OR_ASSIGN, lval, 0
	}
yyrule54: // "||"
	{
		return LOR, lval, 0
	}
yyrule55: // "}"
	{
		return RBRACE, lval, 0
	}
yyrule56: // break
	{
		return BREAK, lval, 0
	}
yyrule57: // case
	{
		return CASE, lval, 0
	}
yyrule58: // chan
	{
		return CHAN, lval, 0
	}
yyrule59: // const
	{
		return CONST, lval, 0
	}
yyrule60: // continue
	{
		return CONTINUE, lval, 0
	}
yyrule61: // default
	{
		return DEFAULT, lval, 0
	}
yyrule62: // defer
	{
		return DEFER, lval, 0
	}
yyrule63: // else
	{
		return ELSE, lval, 0
	}
yyrule64: // fallthrough
	{
		return FALLTHROUGH, lval, 0
	}
yyrule65: // for
	{
		return FOR, lval, 0
	}
yyrule66: // func
	{
		return FUNC, lval, 0
	}
yyrule67: // go
	{
		return GO, lval, 0
	}
yyrule68: // goto
	{
		return GOTO, lval, 0
	}
yyrule69: // if
	{
		return IF, lval, 0
	}
yyrule70: // import
	{
		return IMPORT, lval, 0
	}
yyrule71: // interface
	{
		return INTERFACE, lval, 0
	}
yyrule72: // map
	{
		return MAP, lval, 0
	}
yyrule73: // package
	{
		return PACKAGE, lval, 0
	}
yyrule74: // range
	{
		return RANGE, lval, 0
	}
yyrule75: // return
	{
		return RETURN, lval, 0
	}
yyrule76: // select
	{
		return SELECT, lval, 0
	}
yyrule77: // struct
	{
		return STRUCT, lval, 0
	}
yyrule78: // switch
	{
		return SWITCH, lval, 0
	}
yyrule79: // type
	{
		return GO_TYPE, lval, 0
	}
yyrule80: // var
	{
		return VAR, lval, 0
	}
yyrule81: // {imaginary_ilit}
	{
		return s.int(IMAG)
	}
yyrule82: // {imaginary_lit}
	{
		return s.float(IMAG)
	}
yyrule83: // {int_lit}
	{
		return s.int(INT)
	}
yyrule84: // {float_lit}
	{
		return s.float(FLOAT)
	}
yyrule85: // \"
	{
		s.ssc, s.sc = s.sc, _S1
		goto yystate0
	}
yyrule86: // `
	{
		s.ssc, s.sc = s.sc, _S2
		goto yystate0
	}
yyrule87: // ''
yyrule88: // '(\\.)?[^']*
	{
		return CHAR, string(s.val), 0
	}
yyrule89: // '(\\.)?[^']*'
	{

		s.ssc = s.sc
		if tok, lval, _ = s.str(""); tok != STRING {
			return
		}
		s.i0++
		return CHAR, []rune(lval.(string))[0], 0
	}
yyrule90: // (\\.|[^\\"])*\"
	{
		return s.str("\"")
	}
yyrule91: // ([^`]|\n)*`
	{
		return s.str("`")
	}
yyrule92: // "%{"
	{
		return LCURL, lval, 0
	}
yyrule93: // "%}"
	{
		return RCURL, lval, 0
	}
yyrule94: // %%
	{
		return MARK, lval, 0
	}
yyrule95: // %{errorverbose}
	{
		return ERR_VERBOSE, lval, 0
	}
yyrule96: // %{left}
	{
		return LEFT, lval, 0
	}
yyrule97: // %{nonassoc}
	{
		return NONASSOC, lval, 0
	}
yyrule98: // %{prec}
	{
		return PREC, lval, 0
	}
yyrule99: // %{right}
	{
		return RIGHT, lval, 0
	}
yyrule100: // %{start}
	{
		return START, lval, 0
	}
yyrule101: // %{token}
	{
		return TOKEN, lval, 0
	}
yyrule102: // %{type}
	{
		return TYPE, lval, 0
	}
yyrule103: // %{union}
	{
		return UNION, lval, 0
	}
yyrule104: // ,
	{
		return COMMA, lval, 0
	}
yyrule105: // [a-zA-Z_][a-zA-Z0-9_]*
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
		return IDENT, string(s.val), 0
	}
yyrule106: // [a-zA-Z_][a-zA-Z0-9_.]*
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
		return IDENTIFIER, string(s.val), 0
	}
yyrule107: // {non_ascii}
	{

		s.i = s.i0 - 1
		if rune := s.getRune(false, false); rune < 0 {
			_, sz := utf8.DecodeRune(s.src[s.i:])
			s.i += sz
			s.next()
			s.err("expected unicode lettter, got %U", rune)
			return ILLEGAL, string(-rune), 0
		}
		for rune := rune(1); rune >= 0; rune = s.getRune(true, false) {
		}
		s.next()
		return IDENT, string(s.src[s.i0-1 : s.i-1]), 0
	}
yyrule108: // {non_ascii}
	{

		s.i = s.i0 - 1
		if rune := s.getRune(false, false); rune < 0 {
			_, sz := utf8.DecodeRune(s.src[s.i:])
			s.i += sz
			s.next()
			s.err("expected unicode lettter, got %U", rune)
			return ILLEGAL, string(-rune), 0
		}
		for rune := rune(1); rune >= 0; rune = s.getRune(true, true) {
		}
		s.next()
		return IDENTIFIER, string(s.src[s.i0-1 : s.i-1]), 0
	}
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	s.next()
	return ILLEGAL, string(c0), 0
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

func (s *Scanner) str(pref string) (tok Token, lval interface{}, num int) {
	s.sc = s.ssc
	ss := pref + string(s.val)
	ss, err := strconv.Unquote(ss)
	if err != nil {
		s.err("string literal %q: %v", ss, err)
		return ILLEGAL, ss, 0
	}

	s.i0--
	return STRING, ss, 0
}

func (s *Scanner) int(tk Token) (tok Token, lval interface{}, num int) {
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

func (s *Scanner) float(tk Token) (tok Token, lval interface{}, num int) {
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
