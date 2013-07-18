// Copyright (c) 2013 Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"fmt"
	"path"
	"runtime"
	"testing"
)

func dbg(s string, va ...interface{}) {
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Printf("%s:%d: ", path.Base(fn), fl)
	fmt.Printf(s, va...)
	fmt.Println()
}

type row struct {
	src string
	tok Token
	lit interface{}
}

func testTokens(t *testing.T, yacc bool, table []row) {
	for i, test := range table {
		s := New([]byte(test.src))
		s.Mode(yacc)
		tok, lit := s.Scan()
		if g, e, g2, e2 := tok, test.tok, lit, test.lit; g != e || g2 != e2 {
			t.Errorf("%d: %s(%d) %s(%d)", i, g, g, e, e)
			t.Fatalf("%d: %T(%#v) %T(%#v)", i, g2, g2, e2, e2)
		}
	}
}

func TestGoTokens(t *testing.T) {
	testTokens(t, false, []row{
		{"@", ILLEGAL, "@"}, // 0
		{"%{", REM, nil},
		{"%left", REM, nil},
		{"%%", REM, nil},
		{"%nonassoc", REM, nil},

		{"%prec", REM, nil}, // 5
		{"%}", REM, nil},
		{"%right", REM, nil},
		{"%start", REM, nil},
		{"%token", REM, nil},

		{"%type", REM, nil}, // 10
		{"%union", REM, nil},
		{"", EOF, nil},
		{"//", COMMENT, "//"},
		{"// ", COMMENT, "// "},

		{"/**/ ", COMMENT, "/**/"}, // 15
		{"/***/ ", COMMENT, "/***/"},
		{"/** */ ", COMMENT, "/** */"},
		{"/* **/ ", COMMENT, "/* **/"},
		{"/* * */ ", COMMENT, "/* * */"},

		{"a", IDENT, "a"}, // 20
		{"ab", IDENT, "ab"},
		{"1", INT, uint64(1)},
		{"12", INT, uint64(12)},
		{`""`, STRING, ""},

		{`"1"`, STRING, "1"}, // 25
		{`"12"`, STRING, "12"},
		{"``", STRING, ""},
		{"`1`", STRING, "1"},
		{"`12`", STRING, "12"},

		{"'@'", CHAR, int32('@')}, // 30
		{"a ", IDENT, "a"},
		{"ab ", IDENT, "ab"},
		{"1 ", INT, uint64(1)},
		{"12 ", INT, uint64(12)},

		{`"" `, STRING, ""}, // 35
		{`"1" `, STRING, "1"},
		{`"12" `, STRING, "12"},
		{"`` ", STRING, ""},
		{"`1` ", STRING, "1"},

		{"`12` ", STRING, "12"}, // 40
		{"'@' ", CHAR, int32('@')},
		{" a", IDENT, "a"},
		{" ab", IDENT, "ab"},
		{" 1", INT, uint64(1)},

		{" 12", INT, uint64(12)}, // 45
		{` ""`, STRING, ""},
		{` "1"`, STRING, "1"},
		{` "12"`, STRING, "12"},
		{" ``", STRING, ""},

		{" `1`", STRING, "1"}, // 50
		{" `12`", STRING, "12"},
		{" '@'", CHAR, int32('@')},
		{" a ", IDENT, "a"},
		{" ab ", IDENT, "ab"},

		{" 1 ", INT, uint64(1)}, // 55
		{" 12 ", INT, uint64(12)},
		{` "" `, STRING, ""},
		{` "1" `, STRING, "1"},
		{` "12" `, STRING, "12"},

		{" `` ", STRING, ""}, // 60
		{" `1` ", STRING, "1"},
		{" `12` ", STRING, "12"},
		{" '@' ", CHAR, int32('@')},
		{"f1234567890", IDENT, "f1234567890"},

		{"bár", IDENT, "bár"}, // 65
		{"bára", IDENT, "bára"},
		{"123", INT, uint64(123)},
		{"4e6", FLOAT, 4000000.},
		{"42i", IMAG, 42i},

		{"'@'", CHAR, int32(64)}, // 70
		{`"foo"`, STRING, "foo"},
		{"`foo`", STRING, "foo"},
		{"+", ADD, nil},
		{"-", SUB, nil},

		{"*", MUL, nil}, // 75
		{"/", QUO, nil},
		{"%", REM, nil},
		{"&", AND, nil},
		{"|", OR, nil},

		{"^", XOR, nil}, // 80
		{"<<", SHL, nil},
		{">>", SHR, nil},
		{"&^", AND_NOT, nil},
		{"+=", ADD_ASSIGN, nil},

		{"-=", SUB_ASSIGN, nil}, // 85
		{"*=", MUL_ASSIGN, nil},
		{"/=", QUO_ASSIGN, nil},
		{"%=", REM_ASSIGN, nil},
		{"&=", AND_ASSIGN, nil},

		{"|=", OR_ASSIGN, nil}, // 90
		{"^=", XOR_ASSIGN, nil},
		{"<<=", SHL_ASSIGN, nil},
		{">>=", SHR_ASSIGN, nil},
		{"&^=", AND_NOT_ASSIGN, nil},

		{"&&", LAND, nil}, // 95
		{"||", LOR, nil},
		{"<-", ARROW, nil},
		{"++", INC, nil},
		{"--", DEC, nil},

		{"==", EQL, nil}, // 100
		{"<", LSS, nil},
		{">", GTR, nil},
		{"=", ASSIGN, nil},
		{"!", NOT, nil},

		{"!=", NEQ, nil}, // 105
		{"<=", LEQ, nil},
		{">=", GEQ, nil},
		{":=", DEFINE, nil},
		{"...", ELLIPSIS, nil},

		{"(", LPAREN, nil}, // 110
		{"[", LBRACK, nil},
		{"{", LBRACE, nil},
		{",", COMMA, nil},
		{".", PERIOD, nil},

		{")", RPAREN, nil}, // 115
		{"]", RBRACK, nil},
		{"}", RBRACE, nil},
		{";", SEMICOLON, nil},
		{":", COLON, nil},

		{"break", BREAK, nil}, // 120
		{"case", CASE, nil},
		{"chan", CHAN, nil},
		{"const", CONST, nil},
		{"continue", CONTINUE, nil},

		{"default", DEFAULT, nil}, // 125
		{"defer", DEFER, nil},
		{"else", ELSE, nil},
		{"fallthrough", FALLTHROUGH, nil},
		{"for", FOR, nil},

		{"func", FUNC, nil}, // 130
		{"go", GO, nil},
		{"goto", GOTO, nil},
		{"if", IF, nil},
		{"import", IMPORT, nil},

		{"interface", INTERFACE, nil}, // 135
		{"map", MAP, nil},
		{"package", PACKAGE, nil},
		{"range", RANGE, nil},
		{"return", RETURN, nil},

		{"select", SELECT, nil}, // 140
		{"struct", STRUCT, nil},
		{"switch", SWITCH, nil},
		{"type", GO_TYPE, nil},
		{"var", VAR, nil},
	})
}

func TestYaccTokens(t *testing.T) {
	testTokens(t, true, []row{
		{"@", ILLEGAL, "@"}, // 0
		{"%{", LCURL, nil},
		{"%left", LEFT, nil},
		{"%%", MARK, nil},
		{"%nonassoc", NONASSOC, nil},

		{"%prec", PREC, nil}, // 5
		{"%}", RCURL, nil},
		{"%right", RIGHT, nil},
		{"%start", START, nil},
		{"%token", TOKEN, nil},

		{"%type", TYPE, nil}, // 10
		{"%union", UNION, nil},
		{"", EOF, nil},
		{"//", COMMENT, "//"},
		{"// ", COMMENT, "// "},

		{"/**/ ", COMMENT, "/**/"}, // 15
		{"/***/ ", COMMENT, "/***/"},
		{"/** */ ", COMMENT, "/** */"},
		{"/* **/ ", COMMENT, "/* **/"},
		{"/* * */ ", COMMENT, "/* * */"},

		{"a", IDENTIFIER, "a"}, // 20
		{"ab", IDENTIFIER, "ab"},
		{"1", INT, uint64(1)},
		{"12", INT, uint64(12)},
		{`""`, STRING, ""},

		{`"1"`, STRING, "1"}, // 25
		{`"12"`, STRING, "12"},
		{"``", STRING, ""},
		{"`1`", STRING, "1"},
		{"`12`", STRING, "12"},

		{"'@'", CHAR, int32('@')}, // 30
		{"a ", IDENTIFIER, "a"},
		{"ab ", IDENTIFIER, "ab"},
		{"1 ", INT, uint64(1)},
		{"12 ", INT, uint64(12)},

		{`"" `, STRING, ""}, // 35
		{`"1" `, STRING, "1"},
		{`"12" `, STRING, "12"},
		{"`` ", STRING, ""},
		{"`1` ", STRING, "1"},

		{"`12` ", STRING, "12"}, // 40
		{"'@' ", CHAR, int32('@')},
		{" a", IDENTIFIER, "a"},
		{" ab", IDENTIFIER, "ab"},
		{" 1", INT, uint64(1)},

		{" 12", INT, uint64(12)}, // 45
		{` ""`, STRING, ""},
		{` "1"`, STRING, "1"},
		{` "12"`, STRING, "12"},
		{" ``", STRING, ""},

		{" `1`", STRING, "1"}, // 50
		{" `12`", STRING, "12"},
		{" '@'", CHAR, int32('@')},
		{" a ", IDENTIFIER, "a"},
		{" ab ", IDENTIFIER, "ab"},

		{" 1 ", INT, uint64(1)}, // 55
		{" 12 ", INT, uint64(12)},
		{` "" `, STRING, ""},
		{` "1" `, STRING, "1"},
		{` "12" `, STRING, "12"},

		{" `` ", STRING, ""}, // 60
		{" `1` ", STRING, "1"},
		{" `12` ", STRING, "12"},
		{" '@' ", CHAR, int32('@')},
		{"f1234567890", IDENTIFIER, "f1234567890"},

		{"bár", IDENTIFIER, "b"}, // 65
		{"bára", IDENTIFIER, "b"},
		{"123", INT, uint64(123)},
		{"4e6", INT, uint64(4)},
		{"42i", INT, uint64(42)},

		{"'@'", CHAR, int32(64)}, // 70
		{`"foo"`, STRING, "foo"},
		{"`foo`", STRING, "foo"},
		//TODO{"+", ADD, nil},
		//TODO{"-", SUB, nil},

		//TODO{"*", MUL, nil}, // 75
		//TODO{"/", QUO, nil},
		//TODO{"%", REM, nil},
		//TODO{"&", AND, nil},
		//TODO{"|", OR, nil},

		//TODO{"^", XOR, nil}, // 80
		//TODO{"<<", SHL, nil},
		//TODO{">>", SHR, nil},
		//TODO{"&^", AND_NOT, nil},
		//TODO{"+=", ADD_ASSIGN, nil},

		//TODO{"-=", SUB_ASSIGN, nil}, // 85
		//TODO{"*=", MUL_ASSIGN, nil},
		//TODO{"/=", QUO_ASSIGN, nil},
		//TODO{"%=", REM_ASSIGN, nil},
		//TODO{"&=", AND_ASSIGN, nil},

		//TODO{"|=", OR_ASSIGN, nil}, // 90
		//TODO{"^=", XOR_ASSIGN, nil},
		//TODO{"<<=", SHL_ASSIGN, nil},
		//TODO{">>=", SHR_ASSIGN, nil},
		//TODO{"&^=", AND_NOT_ASSIGN, nil},

		//TODO{"&&", LAND, nil}, // 95
		//TODO{"||", LOR, nil},
		//TODO{"<-", ARROW, nil},
		//TODO{"++", INC, nil},
		//TODO{"--", DEC, nil},

		//TODO{"==", EQL, nil}, // 100
		//TODO{"<", LSS, nil},
		//TODO{">", GTR, nil},
		//TODO{"=", ASSIGN, nil},
		//TODO{"!", NOT, nil},

		//TODO{"!=", NEQ, nil}, // 105
		//TODO{"<=", LEQ, nil},
		//TODO{">=", GEQ, nil},
		//TODO{":=", DEFINE, nil},
		//TODO{"...", ELLIPSIS, nil},

		//TODO{"(", LPAREN, nil}, // 110
		//TODO{"[", LBRACK, nil},
		//TODO{"{", LBRACE, nil},
		//TODO{",", COMMA, nil},
		//TODO{".", PERIOD, nil},

		//TODO{")", RPAREN, nil}, // 115
		//TODO{"]", RBRACK, nil},
		//TODO{"}", RBRACE, nil},
		//TODO{";", SEMICOLON, nil},
		//TODO{":", COLON, nil},

		//TODO{"break", BREAK, nil}, // 120
		//TODO{"case", CASE, nil},
		//TODO{"chan", CHAN, nil},
		//TODO{"const", CONST, nil},
		//TODO{"continue", CONTINUE, nil},

		//TODO{"default", DEFAULT, nil}, // 125
		//TODO{"defer", DEFER, nil},
		//TODO{"else", ELSE, nil},
		//TODO{"fallthrough", FALLTHROUGH, nil},
		//TODO{"for", FOR, nil},

		//TODO{"func", FUNC, nil}, // 130
		//TODO{"go", GO, nil},
		//TODO{"goto", GOTO, nil},
		//TODO{"if", IF, nil},
		//TODO{"import", IMPORT, nil},

		//TODO{"interface", INTERFACE, nil}, // 135
		//TODO{"map", MAP, nil},
		//TODO{"package", PACKAGE, nil},
		//TODO{"range", RANGE, nil},
		//TODO{"return", RETURN, nil},

		//TODO{"select", SELECT, nil}, // 140
		//TODO{"struct", STRUCT, nil},
		//TODO{"switch", SWITCH, nil},
		//TODO{"type", GO_TYPE, nil},
		//TODO{"var", VAR, nil},
	})
}
