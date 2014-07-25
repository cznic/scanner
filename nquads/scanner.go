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
	case c == '#':
		goto yystate16
	case c == '.':
		goto yystate17
	case c == '<':
		goto yystate18
	case c == '@':
		goto yystate29
	case c == '\n' || c == '\r':
		goto yystate4
	case c == '\t' || c == ' ':
		goto yystate3
	case c == '\x00':
		goto yystate2
	case c == '^':
		goto yystate33
	case c == '_':
		goto yystate35
	}

yystate2:
	c = s.next()
	goto yyrule1

yystate3:
	c = s.next()
	switch {
	default:
		goto yyrule2
	case c == '\t' || c == ' ':
		goto yystate3
	}

yystate4:
	c = s.next()
	switch {
	default:
		goto yyrule7
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
	goto yyrule10

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
		goto yyrule3
	case c >= '\x01' && c <= '\t' || c >= '\v' && c <= 'ÿ':
		goto yystate16
	}

yystate17:
	c = s.next()
	goto yyrule4

yystate18:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '!' || c >= '#' && c <= ';' || c == '=' || c >= '?' && c <= '[' || c == ']' || c == '_' || c >= 'a' && c <= 'z' || c >= '~' && c <= 'ÿ':
		goto yystate18
	case c == '>':
		goto yystate19
	case c == '\\':
		goto yystate20
	}

yystate19:
	c = s.next()
	goto yyrule8

yystate20:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == 'U':
		goto yystate21
	case c == 'u':
		goto yystate25
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
		goto yystate27
	}

yystate27:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate28
	}

yystate28:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f':
		goto yystate18
	}

yystate29:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate30
	}

yystate30:
	c = s.next()
	switch {
	default:
		goto yyrule9
	case c == '-':
		goto yystate31
	case c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate30
	}

yystate31:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate32
	}

yystate32:
	c = s.next()
	switch {
	default:
		goto yyrule9
	case c == '-':
		goto yystate31
	case c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		goto yystate32
	}

yystate33:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == '^':
		goto yystate34
	}

yystate34:
	c = s.next()
	goto yyrule5

yystate35:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c == ':':
		goto yystate36
	}

yystate36:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '0' && c <= ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate37
	case c >= 'Â' && c <= 'ß':
		goto yystate38
	case c >= 'à' && c <= 'ï':
		goto yystate39
	case c >= 'ð' && c <= 'ô':
		goto yystate40
	}

yystate37:
	c = s.next()
	switch {
	default:
		goto yyrule6
	case c == '-' || c == '.' || c >= '0' && c <= ':' || c >= 'A' && c <= 'Z' || c == '_' || c >= 'a' && c <= 'z':
		goto yystate37
	case c >= 'Â' && c <= 'ß':
		goto yystate38
	case c >= 'à' && c <= 'ï':
		goto yystate39
	case c >= 'ð' && c <= 'ô':
		goto yystate40
	}

yystate38:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate37
	}

yystate39:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate38
	}

yystate40:
	c = s.next()
	switch {
	default:
		goto yyabort
	case c >= '\u0080' && c <= '¿':
		goto yystate39
	}

yyrule1: // \0
	{
		s.i0++
		return EOF, ""
	}
yyrule2: // [ \t]+

	goto yystate0
yyrule3: // #.*

	goto yystate0
yyrule4: // \.
	{
		return DOT, "."
	}
yyrule5: // "^^"
	{
		return DACCENT, "^^"
	}
yyrule6: // {blank_node_label}
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
	}
yyrule7: // {eol}
	{
		return EOL, ""
	}
yyrule8: // {iriref}
	{

		val, err := strconv.Unquote(`"` + string(s.val) + `"`)
		if err != nil {
			s.err(err.Error())
		}
		return IRIREF, val[1 : len(val)-1]
	}
yyrule9: // {langtag}
	{

		return LANGTAG, string(s.val)
	}
yyrule10: // {string_literal_quote}
	{

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
	}
	panic("unreachable")

	goto yyabort // silence unused label error

yyabort: // no lexem recognized
	s.next()
	return ILLEGAL, string(c0)

}
