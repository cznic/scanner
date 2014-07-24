%{

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
	_ = 0xE000+iota

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
%}


%yyc c
%yyn c = s.next()

langtag	                @[a-zA-Z]+(-[a-zA-Z0-9]+)*
eol	                [\x0D\x0A]+
iriref	                <([^\x01-\x20<>"{}|^`\\]|{uchar})*>
string_literal_quote    \"([^\x22\x5C\x0A\x0D]|{echar}|{uchar})*\"
blank_node_label        _:({pn_chars_u}|[0-9])(({pn_chars}|\.)*{pn_chars})?
uchar                   \\u{hex}{hex}{hex}{hex}|\\U{hex}{hex}{hex}{hex}{hex}{hex}{hex}{hex}
echar                   \\[tbnrf"'\\]
pn_chars_base           [A-Z]|[a-z]|{u}
pn_chars_u              {pn_chars_base}|_|:
pn_chars                {pn_chars_u}|-|[0-9]|{u}
hex                     [0-9]|[A-F]|[a-f]

utf8_h2                 [\xC2-\xDF]
utf8_h3                 [\xE0-\xEF]
utf8_h4                 [\xF0-\xF4]
utf8_cont               [\x80-\xBF]
utf8_2                  {utf8_h2}{utf8_cont}
utf8_3                  {utf8_h3}{utf8_cont}{utf8_cont}
utf8_4                  {utf8_h4}{utf8_cont}{utf8_cont}{utf8_cont}
u                       {utf8_2}|{utf8_3}|{utf8_4}

%%
	s.val = s.val[:0]
	s.i0, s.Line, s.Col, c0 = s.i, s.NLine, s.NCol, c
			
[ \t]+

\0	s.i0++
	return EOF, ""

{blank_node_label}
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
	
	for i := i+n; i < len(s.val); i += n {
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
		s.Col += i-n
		s.err("invalid character %U in BLANK_NODE_LABEL", c)
		return ILLEGAL, ""
	}

	return LABEL, string(s.val[2:])

{eol}
{iriref}
{langtag}
{string_literal_quote}

%%
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
	'A',     'Z',     // 0
	'a',     'z',     // 1
	0x00C0,  0x00D6,  // 2
	0x00D8,  0x00F6,  // 3
	0x00F8,  0x02FF,  // 4
	0x0370,  0x037D,  // 5
	0x037F,  0x1FFF,  // 6
	0x200C,  0x200D,  // 7
	0x2070,  0x218F,  // 8
	0x2C00,  0x2FEF,  // 9
	0x3001,  0xD7FF,  // 10
	0xF900,  0xFDCF,  // 11
	0xFDF0,  0xFFFD,  // 12
	0x10000, 0xEFFFF, // 13, last PN_CHARS_BASE
	'_',     '_',     // 14
	':',     ':',     // 15, last PN_CHARS_U
	'-',     '-',     // 16
	'0',     '9',     // 17
	0x00B7,  0x00B7,  // 18
	0x0300,  0x036F,  // 19
        0x203F,  0x2040,  // 20, last PN_CHARS
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
