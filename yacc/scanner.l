%{

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
	_ = 0xE000+iota

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
	COLON:      "COLON",
	
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
	RETURN: "  RETURN",
	
	SELECT:  "SELECT",
	STRUCT:  "STRUCT",
	SWITCH:  "SWITCH",
	GO_TYPE: "GO_TYPE",
	VAR:      "VAR",

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

	START: "START",
	TOKEN: "TOKEN",
	TYPE:  "TYPE",
	UNION: "UNION",
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
	Col     int     // Starting column of the last scanned token.
	Errors  []error // List of accumulated errors.
	Fname   string  // File name (reported) of the scanned source.
	Line    int     // Starting line of the last scanned token.
	NCol    int     // Starting column (reported) for the next scanned token.
	NLine   int     // Starting line (reported) for the next scanned token.
	c       int
	i       int
	i0      int
	sc      int
	src     []byte
	ssc     int     // saved state condition
	val     []byte
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
%}

int_lit         {decimal_lit}|{octal_lit}|{hex_lit}
decimal_lit     [1-9][0-9]*
octal_lit       0[0-7]*
hex_lit         0[xX][0-9a-fA-F]+

float_lit       {D}"."{D}?{E}?|{D}{E}|"."{D}{E}?
D        	[0-9]+
E        	[eE][-+]?[0-9]+

imaginary_ilit  {D}i
imaginary_lit   {float_lit}i
utf8_h2         [\xC2-\xDF]
utf8_h3         [\xE0-\xEF]
utf8_h4         [\xF0-\xF4]
utf8_cont       [\x80-\xBF]
utf8_2          {utf8_h2}{utf8_cont}
utf8_3          {utf8_h3}{utf8_cont}{utf8_cont}
utf8_4          {utf8_h4}{utf8_cont}{utf8_cont}{utf8_cont}
non_ascii       {utf8_2}|{utf8_3}|{utf8_4}

a		[aA]
b		[bB]
c		[cC]
e		[eE]
f		[fF]
g		[gG]
h		[hH]
i		[iI]
k		[kK]
l		[lL]
n		[nN]
o		[oO]
p		[pP]
r		[rR]
s		[sS]
t		[tT]
u		[uU]
v		[vV]
y		[yY]

errorverbose	{e}{r}{r}{o}{r}-{v}{e}{r}{b}{o}{s}{e}
left		{l}{e}{f}{t}
nonassoc	{n}{o}{n}{a}{s}{s}{o}{c}
prec		{p}{r}{e}{c}
right		{r}{i}{g}{h}{t}
start		{s}{t}{a}{r}{t}
token		{t}{o}{k}{e}{n}
type		{t}{y}{p}{e}
union		{u}{n}{i}{o}{n}

%yyc c
%yyn c = s.next()
%yyt s.sc

%x _GO _S1 _S2

%%
			s.val = s.val[:0]
			s.i0, s.Line, s.Col, c0 = s.i, s.NLine, s.NCol, c
			
<*>\0			s.i0++
			return EOF, lval

<INITIAL,_GO>[ \t\n\r]+

<INITIAL,_GO>\/\*([^*]|\*+[^*/])*\*+\/
			|
<INITIAL,_GO>\/\/.*	return COMMENT, string(s.val)

<_GO>"!"		return NOT, lval
<_GO>"!="		return NEQ, lval
<_GO>"%"		return REM, lval
<_GO>"%="		return REM_ASSIGN, lval
<_GO>"&"		return AND, lval
<_GO>"&&"		return LAND, lval
<_GO>"&="		return AND_ASSIGN, lval
<_GO>"&^"		return AND_NOT, lval
<_GO>"&^="		return AND_NOT_ASSIGN, lval
<_GO>"("		return LPAREN, lval
<_GO>")"		return RPAREN, lval
<_GO>"*"		return MUL, lval
<_GO>"*="		return MUL_ASSIGN, lval
<_GO>"+"		return ADD, lval
<_GO>"++"		return INC, lval
<_GO>"+="		return ADD_ASSIGN, lval
<_GO>","		return COMMA, lval
<_GO>"-"		return SUB, lval
<_GO>"--"		return DEC, lval
<_GO>"-="		return SUB_ASSIGN, lval
<_GO>"."		return PERIOD, lval
<_GO>"..."		return ELLIPSIS, lval
<_GO>"/"		return QUO, lval
<_GO>"/="		return QUO_ASSIGN, lval
<_GO>":"		return COLON, lval
<_GO>":="		return DEFINE, lval
<_GO>";"		return SEMICOLON, lval
<_GO>"<"		return LSS, lval
<_GO>"<-"		return ARROW, lval
<_GO>"<<"		return SHL, lval
<_GO>"<<="		return SHL_ASSIGN, lval
<_GO>"<="		return LEQ, lval
<_GO>"="		return ASSIGN, lval
<_GO>"=="		return EQL, lval
<_GO>">"		return GTR, lval
<_GO>">="		return GEQ, lval
<_GO>">>"		return SHR, lval
<_GO>">>="		return SHR_ASSIGN, lval
<_GO>"["		return LBRACK, lval
<_GO>"]"		return RBRACK, lval
<_GO>"^"		return XOR, lval
<_GO>"^="		return XOR_ASSIGN, lval
<_GO>"{"		return LBRACE, lval
<_GO>"|"		return OR, lval
<_GO>"|="		return OR_ASSIGN, lval
<_GO>"||"		return LOR, lval
<_GO>"}"			return RBRACE, lval

<_GO>break		return BREAK, lval
<_GO>case		return CASE, lval
<_GO>chan		return CHAN, lval
<_GO>const		return CONST, lval
<_GO>continue		return CONTINUE, lval
<_GO>default		return DEFAULT, lval
<_GO>defer		return DEFER, lval
<_GO>else		return ELSE, lval
<_GO>fallthrough	return FALLTHROUGH, lval
<_GO>for		return FOR, lval
<_GO>func		return FUNC, lval
<_GO>go			return GO, lval
<_GO>goto		return GOTO, lval
<_GO>if			return IF, lval
<_GO>import		return IMPORT, lval
<_GO>interface		return INTERFACE, lval
<_GO>map		return MAP, lval
<_GO>package		return PACKAGE, lval
<_GO>range		return RANGE, lval
<_GO>return		return RETURN, lval
<_GO>select		return SELECT, lval
<_GO>struct		return STRUCT, lval
<_GO>switch		return SWITCH, lval
<_GO>type		return GO_TYPE, lval
<_GO>var		return VAR, lval

<_GO>{imaginary_ilit}	return s.int(IMAG)
<_GO>{imaginary_lit}	return s.float(IMAG)
<INITIAL,_GO>{int_lit}	return s.int(INT)
<_GO>{float_lit}	return s.float(FLOAT)

<INITIAL,_GO>\"		s.ssc, s.sc = s.sc, _S1
<INITIAL,_GO>`		s.ssc, s.sc = s.sc, _S2

<INITIAL,_GO>''			|
<INITIAL,_GO>'(\\.)?[^']*	return CHAR, string(s.val)

<INITIAL,_GO>'(\\.)?[^']*'	
			s.ssc = s.sc
			if tok, lval = s.str(""); tok != STRING {
				return
			}
			s.i0++
			return CHAR, int32(lval.(string)[0])

<_S1>(\\.|[^\\"])*\"	return s.str("\"")
<_S2>([^`]|\n)*`	return s.str("`")

"%{"			return LCURL, lval
"%}"			return RCURL, lval
%%			return MARK, lval
%{errorverbose}		return ERR_VERBOSE, lval
%{left}			return LEFT, lval
%{nonassoc}		return NONASSOC, lval
%{prec}			return PREC, lval
%{right}		return RIGHT, lval
%{start}		return START, lval
%{token}		return TOKEN, lval
%{type}			return TYPE, lval
%{union}		return UNION, lval

<_GO>[a-zA-Z_][a-zA-Z0-9_]*
	if c >= '\xC2' && c <= '\xF4' {
		s.i--
		s.NCol--
		for rune := rune(1) ; rune >= 0; rune = s.getRune(true, false) {
		tok, lval = IDENT, string(s.src[s.i0-1:s.i])
		}
		s.next()
		return
	}

	return IDENT, string(s.val)

<INITIAL>[a-zA-Z_][a-zA-Z0-9_.]*
	if c >= '\xC2' && c <= '\xF4' {
		s.i--
		s.NCol--
		for rune := rune(1) ; rune >= 0; rune = s.getRune(true, true) {
		tok, lval = IDENTIFIER, string(s.src[s.i0-1:s.i])
		}
		s.next()
		return
	}

	return IDENTIFIER, string(s.val)

<_GO>{non_ascii}
	s.i = s.i0-1
	if rune := s.getRune(false, false); rune < 0 {
		_, sz := utf8.DecodeRune(s.src[s.i:])
		s.i += sz
		s.next()
		s.err("expected unicode lettter, got %U", rune)
		return ILLEGAL, string(-rune)
	}

	for rune := rune(1) ; rune >= 0; rune = s.getRune(true, false) {
	}
	s.next()
	return IDENT, string(s.src[s.i0-1:s.i-1])

<INITIAL>{non_ascii}
	s.i = s.i0-1
	if rune := s.getRune(false, false); rune < 0 {
		_, sz := utf8.DecodeRune(s.src[s.i:])
		s.i += sz
		s.next()
		s.err("expected unicode lettter, got %U", rune)
		return ILLEGAL, string(-rune)
	}

	for rune := rune(1) ; rune >= 0; rune = s.getRune(true, true) {
	}
	s.next()
	return IDENTIFIER, string(s.src[s.i0-1:s.i-1])

%%
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
