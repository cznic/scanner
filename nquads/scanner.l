%{
// Copyright (c) 2014 Go Authors. All rights reserved.  Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.
// 
// CAUTION: If this file is a Go source file (*.go), it was generated
// automatically by '$ golex' from a *.l file - DO NOT EDIT in that case!

package scanner

import (
	"strconv"
)

// Scan scans the next token and returns the token and its value if applicable.
// The source end is indicated by EOF.
// 
// If the returned token is ILLEGAL, the literal string is the offending
// character or number/string/char literal.
func (s *Scanner) Scan() (tok Token, st string) {
	//defer func() { dbg("%s:%d:%d %v %q :%d:%d", s.Fname, s.Line, s.Col, tok, st, s.NLine, s.NCol) }()
	c0, c := s.c, s.c
%}


%yyc c
%yyn c = s.next()

langtag	                @[a-zA-Z]+(-[a-zA-Z0-9]+)*
eol	                [\x0D\x0A]+
iriref	                <([^\x01-\x20<>"{}|^`\\]|{uchar})*>
string_literal_quote    \"([^\x22\x5C\x0A\x0D]|{echar}|{uchar})*\"
blank_node_label        _:({pn_chars_u}|[0-9])({pn_chars}|\.)*
uchar                   \\u{hex}{hex}{hex}{hex}|\\U{hex}{hex}{hex}{hex}{hex}{hex}{hex}{hex}
echar                   \\[tbnrf"'\\]
pn_chars_base           [A-Z]|[a-z]|{na}
pn_chars_u              {pn_chars_base}|_|:
pn_chars                {pn_chars_u}|-|[0-9]|{na}
hex                     [0-9]|[A-F]|[a-f]
na                      [\xC2-\xDF]{n}|[\xE0-\xEF]{n}{n}|[\xF0-\xF4]{n}{n}{n}  // non ascii
n                       [\x80-\xBF]

%%
	s.val = s.val[:0]
	s.i0, s.Line, s.Col, c0 = s.i, s.NLine, s.NCol, c
			
\0	s.i0++
	return EOF, ""

[ \t]+
#.*

\.	return DOT, "."

"^^"	return DACCENT, "^^"

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

	if c == '.' || !checkPnChars(c) {
		for {
			s.val = s.val[:len(s.val)-1]
			s.back()
			if s.val[len(s.val)-1] != '.' {
				s.i++
				break
			}
		}
	}

	return LABEL, string(s.val[2:])	

{eol}	return EOL, ""

{iriref}
	val, err := strconv.Unquote(`"`+string(s.val)+`"`)
	if err != nil {
		s.err(err.Error())
	}
	return IRIREF, val[1:len(val)-1]

{langtag}
	return LANGTAG, string(s.val)

{string_literal_quote}
	// \' needs special preprocessing.
	var c, prev byte
	v := s.val
	for i := 1; i < len(v)-2; i, prev = i+1, c {
		if c = v[i]; prev != '\\' && c == '\\' && v[i+1] == '\'' {
			v = append(v[:i], v[i+1:]...)
		}
	}
	val, err := strconv.Unquote(string(v))
	if err != nil {
		s.err(err.Error())
	}
	return STRING, val

%%
	s.next()
	return ILLEGAL, string(c0)

}
