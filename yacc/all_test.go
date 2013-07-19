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
		{"+", ILLEGAL, "+"},
		{"-", ILLEGAL, "-"},

		{"*", ILLEGAL, "*"}, // 75
		{"/", ILLEGAL, "/"},
		{"%", ILLEGAL, "%"},
		{"&", ILLEGAL, "&"},
		{"|", ILLEGAL, "|"},

		{"^", ILLEGAL, "^"}, // 80
		{"<<", ILLEGAL, "<"},
		{">>", ILLEGAL, ">"},
		{"&^", ILLEGAL, "&"},
		{"+=", ILLEGAL, "+"},

		{"-=", ILLEGAL, "-"}, // 85
		{"*=", ILLEGAL, "*"},
		{"/=", ILLEGAL, "/"},
		{"%=", ILLEGAL, "%"},
		{"&=", ILLEGAL, "&"},

		{"|=", ILLEGAL, "|"}, // 90
		{"^=", ILLEGAL, "^"},
		{"<<=", ILLEGAL, "<"},
		{">>=", ILLEGAL, ">"},
		{"&^=", ILLEGAL, "&"},

		{"&&", ILLEGAL, "&"}, // 95
		{"||", ILLEGAL, "|"},
		{"<-", ILLEGAL, "<"},
		{"++", ILLEGAL, "+"},
		{"--", ILLEGAL, "-"},

		{"==", ILLEGAL, "="}, // 100
		{"<", ILLEGAL, "<"},
		{">", ILLEGAL, ">"},
		{"=", ILLEGAL, "="},
		{"!", ILLEGAL, "!"},

		{"!=", ILLEGAL, "!"}, // 105
		{"<=", ILLEGAL, "<"},
		{">=", ILLEGAL, ">"},
		{":=", ILLEGAL, ":"},
		{"...", ILLEGAL, "."},

		{"(", ILLEGAL, "("}, // 110
		{"[", ILLEGAL, "["},
		{"{", ILLEGAL, "{"},
		{",", ILLEGAL, ","},
		{".", ILLEGAL, "."},

		{")", ILLEGAL, ")"}, // 115
		{"]", ILLEGAL, "]"},
		{"}", ILLEGAL, "}"},
		{";", ILLEGAL, ";"},
		{":", ILLEGAL, ":"},

		{"break", IDENTIFIER, "break"}, // 120
		{"case", IDENTIFIER, "case"},
		{"chan", IDENTIFIER, "chan"},
		{"const", IDENTIFIER, "const"},
		{"continue", IDENTIFIER, "continue"},

		{"default", IDENTIFIER, "default"}, // 125
		{"defer", IDENTIFIER, "defer"},
		{"else", IDENTIFIER, "else"},
		{"fallthrough", IDENTIFIER, "fallthrough"},
		{"for", IDENTIFIER, "for"},

		{"func", IDENTIFIER, "func"}, // 130
		{"go", IDENTIFIER, "go"},
		{"goto", IDENTIFIER, "goto"},
		{"if", IDENTIFIER, "if"},
		{"import", IDENTIFIER, "import"},

		{"interface", IDENTIFIER, "interface"}, // 135
		{"map", IDENTIFIER, "map"},
		{"package", IDENTIFIER, "package"},
		{"range", IDENTIFIER, "range"},
		{"return", IDENTIFIER, "return"},

		{"select", IDENTIFIER, "select"}, // 140
		{"struct", IDENTIFIER, "struct"},
		{"switch", IDENTIFIER, "switch"},
		{"type", IDENTIFIER, "type"},
		{"var", IDENTIFIER, "var"},

		// ----

		{"a.foo", IDENTIFIER, "a.foo"}, // 145
		{"b.fooára", IDENTIFIER, "b.foo"},
		{"ab.foo", IDENTIFIER, "ab.foo"},
		{"a.foo ", IDENTIFIER, "a.foo"},
		{"ab.foo ", IDENTIFIER, "ab.foo"},
		{" a.foo", IDENTIFIER, "a.foo"},

		{" ab.foo", IDENTIFIER, "ab.foo"}, // 150
		{" a.foo ", IDENTIFIER, "a.foo"},
		{" ab.foo ", IDENTIFIER, "ab.foo"},
		{"f1234567890.foo", IDENTIFIER, "f1234567890.foo"},
		{"b.fooár", IDENTIFIER, "b.foo"},

		// --

		{"a:", C_IDENTIFIER, "a"}, // 155
		{"a: ", C_IDENTIFIER, "a"},
		{"a :", C_IDENTIFIER, "a"},
		{"a : ", C_IDENTIFIER, "a"},
		{" a:", C_IDENTIFIER, "a"},

		{" a: ", C_IDENTIFIER, "a"}, // 160
		{" a :", C_IDENTIFIER, "a"},
		{" a : ", C_IDENTIFIER, "a"},
		{"ab:", C_IDENTIFIER, "ab"},
		{"ab: ", C_IDENTIFIER, "ab"},

		{"ab :", C_IDENTIFIER, "ab"}, // 165
		{"ab : ", C_IDENTIFIER, "ab"},
		{" ab:", C_IDENTIFIER, "ab"},
		{" ab: ", C_IDENTIFIER, "ab"},
		{" ab :", C_IDENTIFIER, "ab"},

		{" ab : ", C_IDENTIFIER, "ab"}, // 170
		{"a.b:", C_IDENTIFIER, "a.b"},
		{"a.b: ", C_IDENTIFIER, "a.b"},
		{"a.b :", C_IDENTIFIER, "a.b"},
		{"a.b : ", C_IDENTIFIER, "a.b"},

		{" a.b:", C_IDENTIFIER, "a.b"}, // 175
		{" a.b: ", C_IDENTIFIER, "a.b"},
		{" a.b :", C_IDENTIFIER, "a.b"},
		{" a.b : ", C_IDENTIFIER, "a.b"},
	})
}

func TestBug(t *testing.T) {
	tab := []struct {
		src  string
		toks []Token
	}{
		{`%left`, []Token{LEFT}},
		{`%left %left`, []Token{LEFT, LEFT}},
		{`%left 'a' %left`, []Token{LEFT, CHAR, LEFT}},
		{`foo`, []Token{IDENTIFIER}},
		{`foo bar`, []Token{IDENTIFIER, IDENTIFIER}},
		{`foo bar baz`, []Token{IDENTIFIER, IDENTIFIER, IDENTIFIER}},
		{`%token <ival> DREG VREG`, []Token{TOKEN, ILLEGAL, IDENTIFIER, ILLEGAL, IDENTIFIER, IDENTIFIER}},
	}

	for i, test := range tab {
		s := New([]byte(test.src))
		s.Mode(true)
		for j, etok := range test.toks {
			tok, _ := s.Scan()
			if g, e := tok, etok; g != e {
				t.Errorf("%d.%d: %s(%d) %s(%d)", i, j, g, g, e, e)
			}
		}
	}
}
