%{

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
	"strconv"
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

%yyc c
%yyn c = s.next()
%yyt s.sc

%x S1 S2

%%
			s.val = s.val[:0]
			s.i0, s.Line, s.Col, c0 = s.i, s.NLine, s.NCol, c
			
<*>\0			s.i0++
			return token.EOF, lval

[ \t\r]+

\/\*([^*]|\*+[^*/])*\*+\/	|
\/\/.*			return token.COMMENT, string(s.val)

"!"			return token.NOT, lval
"!="			return token.NEQ, lval
"%"			return token.REM, lval
"%="			return token.REM_ASSIGN, lval
"&"			return token.AND, lval
"&&"			return token.LAND, lval
"&="			return token.AND_ASSIGN, lval
"&^"			return token.AND_NOT, lval
"&^="			return token.AND_NOT_ASSIGN, lval
"("			return token.LPAREN, lval
")"			return token.RPAREN, lval
"*"			return token.MUL, lval
"*="			return token.MUL_ASSIGN, lval
"+"			return token.ADD, lval
"++"			return token.INC, lval
"+="			return token.ADD_ASSIGN, lval
","			return token.COMMA, lval
"-"			return token.SUB, lval
"--"			return token.DEC, lval
"-="			return token.SUB_ASSIGN, lval
"."			return token.PERIOD, lval
"..."			return token.ELLIPSIS, lval
"/"			return token.QUO, lval
"/="			return token.QUO_ASSIGN, lval
":"			return token.COLON, lval
":="			return token.DEFINE, lval
";"			return token.SEMICOLON, lval
"<"			return token.LSS, lval
"<-"			return token.ARROW, lval
"<<"			return token.SHL, lval
"<<="			return token.SHL_ASSIGN, lval
"<="			return token.LEQ, lval
"="			return token.ASSIGN, lval
"=="			return token.EQL, lval
">"			return token.GTR, lval
">="			return token.GEQ, lval
">>"			return token.SHR, lval
">>="			return token.SHR_ASSIGN, lval
"["			return token.LBRACK, lval
"]"			return token.RBRACK, lval
"^"			return token.XOR, lval
"^="			return token.XOR_ASSIGN, lval
"{"			return token.LBRACE, lval
"|"			return token.OR, lval
"|="			return token.OR_ASSIGN, lval
"||"			return token.LOR, lval
"}"			return token.RBRACE, lval

break			return token.BREAK, lval
case			return token.CASE, lval
chan			return token.CHAN, lval
const			return token.CONST, lval
continue		return token.CONTINUE, lval
default			return token.DEFAULT, lval
defer			return token.DEFER, lval
else			return token.ELSE, lval
fallthrough		return token.FALLTHROUGH, lval
for			return token.FOR, lval
func			return token.FUNC, lval
go			return token.GO, lval
goto			return token.GOTO, lval
if			return token.IF, lval
import			return token.IMPORT, lval
interface		return token.INTERFACE, lval
map			return token.MAP, lval
package			return token.PACKAGE, lval
range			return token.RANGE, lval
return			return token.RETURN, lval
select			return token.SELECT, lval
struct			return token.STRUCT, lval
switch			return token.SWITCH, lval
type			return token.TYPE, lval
var			return token.VAR, lval

{imaginary_ilit}	return s.int(token.IMAG)
{imaginary_lit}		return s.float(token.IMAG)
{int_lit}		return s.int(token.INT)
{float_lit}		return s.float(token.FLOAT)

\"			s.sc = S1
`			s.sc = S2

''			return token.CHAR, string(s.val)
'(\\.)?[^']*		return token.CHAR, string(s.val)

'(\\.)?[^']*'		if tok, lval = s.str(""); tok != token.STRING {
				tok = token.CHAR
				return
			}
			s.i0++
			ss := lval.(string)
			if ss != "" {
				return token.CHAR, int32(ss[0])
			}

			return token.CHAR, string(s.val)

<S1>(\\.|[^\\"])*\"	return s.str("\"")
<S2>([^`]|\n)*`		return s.str("`")

[a-zA-Z_][a-zA-Z_0-9]*
	if c >= '\xC2' && c <= '\xF4' {
		s.i--
		s.NCol--
		for rune := rune(1) ; rune >= 0; rune = s.getRune(true) {
		}
		tok, lval = token.IDENT, string(s.src[s.i0-1:s.i])
		s.next()
		return
	}
	return token.IDENT, string(s.val)

{non_ascii}
	s.i = s.i0-1
	if rune := s.getRune(false); rune < 0 {
		_, sz := utf8.DecodeRune(s.src[s.i:])
		s.i += sz
		s.next()
		s.err("expected unicode lettter, got %U", rune)
		return token.ILLEGAL, string(-rune)
	}

	for rune := rune(1) ; rune >= 0; rune = s.getRune(true) {
	}
	s.next()
	return token.IDENT, string(s.src[s.i0-1:s.i-1])

%%
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
	ss0 := pref + string(s.val)
	ss, err := strconv.Unquote(ss0)
	if err != nil {
		s.err("string literal %q: %v", ss, err)
		s.i0--
		return token.STRING, ss0
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
