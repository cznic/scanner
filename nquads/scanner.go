// Copyright (c) 2014 Go Authors. All rights reserved.  Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.
//
// CAUTION: If this file is a Go source file (*.go), it was generated
// automatically by '$ golex' from a *.l file - DO NOT EDIT in that case!

// Package scanner implements a scanner for N-Quads[0] source text. It takes a
// []byte as source which can then be tokenized through repeated calls to the
// Scan method.
//
// Links
//
// Referenced from above:
//
// [0]: http://www.w3.org/TR/n-quads/
package scanner

import (
	"errors"
	"fmt"
	"math"
	"unicode"
)

// Productions for terminals
//
// LANGTAG		::=	'@' [a-zA-Z]+ ('-' [a-zA-Z0-9]+)*
// EOL			::=	[#xD#xA]+
// IRIREF		::=	'<' ([^#x00-#x20<>"{}|^`\] | UCHAR)* '>'
// STRING_LITERAL_QUOTE	::=	'"' ([^#x22#x5C#xA#xD] | ECHAR | UCHAR)* '"'
// BLANK_NODE_LABEL	::=	'_:' (PN_CHARS_U | [0-9]) ((PN_CHARS | '.')* PN_CHARS)?
// UCHAR		::=	'\u' HEX HEX HEX HEX | '\U' HEX HEX HEX HEX HEX HEX HEX HEX
// ECHAR		::=	'\' [tbnrf"'\]
// PN_CHARS_BASE	::=	[A-Z] | [a-z] | [#x00C0-#x00D6] | [#x00D8-#x00F6]
// 	| [#x00F8-#x02FF] | [#x0370-#x037D] | [#x037F-#x1FFF] | [#x200C-#x200D]
// 	| [#x2070-#x218F] | [#x2C00-#x2FEF] | [#x3001-#xD7FF] | [#xF900-#xFDCF]
// 	| [#xFDF0-#xFFFD] | [#x10000-#xEFFFF]
// PN_CHARS_U		::=	PN_CHARS_BASE | '_' | ':'
// PN_CHARS		::=	PN_CHARS_U | '-' | [0-9] | #x00B7 | [#x0300-#x036F] | [#x203F-#x2040]
// HEX			::=	[0-9] | [A-F] | [a-f]

// Token is the type of the token identifier returned by Scan().
type Token int

// Values of Token.
//  [0]: http://www.w3.org/TR/n-quads/#grammar-production-BLANK_NODE_LABEL
//  [1]: http://www.w3.org/TR/n-quads/#grammar-production-EOL
//  [2]: http://www.w3.org/TR/n-quads/#grammar-production-IRIREF
//  [3]: http://www.w3.org/TR/n-quads/#grammar-production-LANGTAG
//  [4]: http://www.w3.org/TR/n-quads/#grammar-production-STRING_LITERAL_QUOTE
const (
	_ = 0xE000 + iota

	// ------------------------------------------- N-Quads terminals

	// Special tokens
	ILLEGAL Token = iota
	EOF

	LABEL   // [0]
	EOL     // [1]
	IRIREF  // [2]
	LANGTAG // [3]
	STRING  // [4]
)

var ts = map[Token]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	LABEL:   "LABEL",
	EOL:     "EOL",
	IRIREF:  "IRIREF",
	LANGTAG: "LANGTAG",
	STRING:  "STRING",
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
		switch s.c {
		case '\n':
			s.NLine++
			s.NCol = 0
			if s.i == len(s.src) {
				s.NCol = 1
			}
		default:
			s.NCol++
		}
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

// Scan scans the next token and returns the token and its value if applicable.
// The source end is indicated by EOF.
//
// If the returned token is ILLEGAL, the literal string is the offending
// character or number/string/char literal.
func (s *Scanner) Scan() (Token, string) {
	c0, c := s.c, s.c

yystate0:

	s.val = s.val[:0]
	s.i0, s.Line, s.Col, c0 = s.i, s.NLine, s.NCol, c

	goto yystart1

	goto yystate1 // silence unused label error
yystate1:
	c = s.next()
yystart1:
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate5
	case c == '<':
		goto yystate16
	case c == '@':
		goto yystate27
	case c == '\n' || c == '\r':
		goto yystate4
	case c == '\t' || c == ' ':
		goto yystate3
	case c == '\x00':
		goto yystate2
	case c == '_':
		goto yystate31
	}

yystate2:
	c = s.next()
	goto yyrule2

yystate3:
	c = s.next()
	switch {
	default:
		goto yyrule1
	case c == '\t' || c == ' ':
		goto yystate3
	}

yystate4:
	c = s.next()
	switch {
	default:
		goto yyrule4
	case c == '\n' || c == '\r':
		goto yystate4
	}

yystate5:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '"':
		goto yystate6
	case c == '\\':
		goto yystate7
	case c >= '\x01' && c <= '\t' || c == '\v' || c == '\f' || c >= '\x0e' && c <= '!' || c >= '#' && c <= '[' || c >= ']' && c <= 'ÿ':
		goto yystate5
	}

yystate6:
	c = s.next()
	goto yyrule7

yystate7:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '"' || c == '\'' || c == '\\' || c == 'b' || c == 'f' || c == 'n' || c == 'r' || c == 't':
		goto yystate5
	case c == 'U':
		goto yystate8
	case c == 'u':
		goto yystate12
	}

yystate8:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate9
	}

yystate9:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate10
	}

yystate10:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate11
	}

yystate11:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate12
	}

yystate12:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate13
	}

yystate13:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate14
	}

yystate14:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate15
	}

yystate15:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate5
	}

yystate16:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '!' || c >= '#' && c <= ';' || c == '=' || c >= '?' && c <= '[' || c == ']' || c == '_' || c >= 'a' && c <= 'z' || c >= '~' && c <= 'ÿ':
		goto yystate16
	case c == '>':
		goto yystate17
	case c == '\\':
		goto yystate18
	}

yystate17:
	c = s.next()
	goto yyrule5

yystate18:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'U':
		goto yystate19
	case c == 'u':
		goto yystate23
	}

yystate19:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate20
	}

yystate20:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate21
	}

yystate21:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate22
	}

yystate22:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate23
	}

yystate23:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate24
	}

yystate24:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate25
	}

yystate25:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate26
	}

yystate26:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate16
	}

yystate27:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate28
	}

yystate28:
	c = s.next()
	switch {
	default:
		goto yyrule6
	case c == '-':
		goto yystate29
	case c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate28
	}

yystate29:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate30
	}

yystate30:
	c = s.next()
	switch {
	default:
		goto yyrule6
	case c == '-':
		goto yystate29
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate30
	}

yystate31:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == ':':
		goto yystate32
	}

yystate32:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate33
	case c >= 'Â' && c <= 'ß':
		goto yystate35
	case c >= 'à' && c <= 'ï':
		goto yystate36
	case c >= 'ð' && c <= 'ô':
		goto yystate37
	}

yystate33:
	c = s.next()
	switch {
	default:
		goto yyrule3
	case c == '-' || c >= '0' && c <= ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate33
	case c == '.':
		goto yystate34
	case c >= 'Â' && c <= 'ß':
		goto yystate35
	case c >= 'à' && c <= 'ï':
		goto yystate36
	case c >= 'ð' && c <= 'ô':
		goto yystate37
	}

yystate34:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '-' || c >= '0' && c <= ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate33
	case c == '.':
		goto yystate34
	case c >= 'Â' && c <= 'ß':
		goto yystate35
	case c >= 'à' && c <= 'ï':
		goto yystate36
	case c >= 'ð' && c <= 'ô':
		goto yystate37
	}

yystate35:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate33
	}

yystate36:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate35
	}

yystate37:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate36
	}

yyrule1: // [ \t]+

	goto yystate0
yyrule2: // \0
	{
		s.i0++
		return EOF, ""
	}
yyrule3: // {blank_node_label}
	{

		i := 2
		c, n := decodeRune(s.val[i:])
		switch {
		case c >= '0' && c <= '9', checkPnCharsU(c):
			// ok
		default:
			s.Col += i
			s.err("invalid character %U in BLANK_NODE_LABEL", c)
			return ILLEGAL, ""
		}

		for i := i + n; i < len(s.val); i += n {
			c, n = decodeRune(s.val[i:])
			switch {
			case c == '.', checkPnChars(c):
				// ok
			default:
				s.Col += i
				s.err("invalid character %U in BLANK_NODE_LABEL", c)
				return ILLEGAL, ""
			}
		}
		if c == '.' {
			s.Col += i - n
			s.err("invalid character %U in BLANK_NODE_LABEL", c)
			return ILLEGAL, ""
		}
		return LABEL, string(s.val[2:])
	}
yyrule4: // {eol}

	goto yystate0
yyrule5: // {iriref}

	goto yystate0
yyrule6: // {langtag}

	goto yystate0
yyrule7: // {string_literal_quote}

	goto yystate0
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	s.next()
	return ILLEGAL, string(c0)

}

const (
	RuneError = math.MaxInt32

	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000

	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1
)

func decodeRune(s []byte) (r rune, size int) {
	n := len(s)
	if n < 1 {
		return RuneError, 0
	}
	c0 := s[0]

	// 1-byte, 7-bit sequence?
	if c0 < tx {
		return rune(c0), 1
	}

	// unexpected continuation byte?
	if c0 < t2 {
		return RuneError, 1
	}

	// need first continuation byte
	if n < 2 {
		return RuneError, 1
	}
	c1 := s[1]
	if c1 < tx || t2 <= c1 {
		return RuneError, 1
	}

	// 2-byte, 11-bit sequence?
	if c0 < t3 {
		r = rune(c0&mask2)<<6 | rune(c1&maskx)
		if r <= rune1Max {
			return RuneError, 1
		}
		return r, 2
	}

	// need second continuation byte
	if n < 3 {
		return RuneError, 1
	}
	c2 := s[2]
	if c2 < tx || t2 <= c2 {
		return RuneError, 1
	}

	// 3-byte, 16-bit sequence?
	if c0 < t4 {
		r = rune(c0&mask3)<<12 | rune(c1&maskx)<<6 | rune(c2&maskx)
		if r <= rune2Max {
			return RuneError, 1
		}
		return r, 3
	}

	// need third continuation byte
	if n < 4 {
		return RuneError, 1
	}
	c3 := s[3]
	if c3 < tx || t2 <= c3 {
		return RuneError, 1
	}

	// 4-byte, 21-bit sequence?
	if c0 < t5 {
		r = rune(c0&mask4)<<18 | rune(c1&maskx)<<12 | rune(c2&maskx)<<6 | rune(c3&maskx)
		if r <= rune3Max || unicode.MaxRune < r {
			return RuneError, 1
		}
		return r, 4
	}

	// error
	return RuneError, 1
}

var tab = []rune{
	'A', 'Z', // 0
	'a', 'z', // 1
	0x00C0, 0x00D6, // 2
	0x00D8, 0x00F6, // 3
	0x00F8, 0x02FF, // 4
	0x0370, 0x037D, // 5
	0x037F, 0x1FFF, // 6
	0x200C, 0x200D, // 7
	0x2070, 0x218F, // 8
	0x2C00, 0x2FEF, // 9
	0x3001, 0xD7FF, // 10
	0xF900, 0xFDCF, // 11
	0xFDF0, 0xFFFD, // 12
	0x10000, 0xEFFFF, // 13, last PN_CHARS_BASE
	'_', '_', // 14
	':', ':', // 15, last PN_CHARS_U
	'-', '-', // 16
	'0', '9', // 17
	0x00B7, 0x00B7, // 18
	0x0300, 0x036F, // 19
	0x203F, 0x2040, // 20, last PN_CHARS
}

func check(r rune, tab []rune) bool {
	for i := 0; i < len(tab); i += 2 {
		if r >= tab[i] && r <= tab[i+1] {
			return true
		}
	}
	return false
}

func checkPnCharsU(r rune) bool {
	return check(r, tab[:2*16])
}

func checkPnChars(r rune) bool {
	return check(r, tab)
}
