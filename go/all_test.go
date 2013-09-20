// Copyright (c) 2013 Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"fmt"
	goscanner "go/scanner"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

//TODO test comparative errorlist generation

func dbg(s string, va ...interface{}) {
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Printf("%s:%d: ", path.Base(fn), fl)
	fmt.Printf(s, va...)
	fmt.Println()
}

var (
	std       = filepath.Join(runtime.GOROOT(), "src")
	tests     = filepath.Join(runtime.GOROOT(), "test")
	whitelist = []string{
		filepath.Join("fixedbugs", "bug169.go"), // go/scanner doesn't return the last \n
	}
)

type row struct {
	src string
	tok token.Token
	lit interface{}
}

func testTokens(t *testing.T, yacc bool, table []row) {
	for i, test := range table {
		s := New([]byte(test.src))
		tok, lit := s.Scan()
		if g, e := tok, test.tok; g != e {
			t.Error(i, g, e)
		}

		if g, e := lit, test.lit; g != e {
			t.Errorf("%d: %T(%#v) %T(%#v)", i, g, g, e, e)
		}
	}
}

func TestGoTokens(t *testing.T) {
	testTokens(t, false, []row{
		{"@", token.ILLEGAL, "@"}, // 0
		{"", token.EOF, nil},
		{"//", token.COMMENT, "//"},
		{"// ", token.COMMENT, "// "},
		{"/**/ ", token.COMMENT, "/**/"},

		{"/***/ ", token.COMMENT, "/***/"}, // 5
		{"/** */ ", token.COMMENT, "/** */"},
		{"/* **/ ", token.COMMENT, "/* **/"},
		{"/* * */ ", token.COMMENT, "/* * */"},
		{"a", token.IDENT, "a"},

		{"ab", token.IDENT, "ab"}, // 10
		{"1", token.INT, uint64(1)},
		{"12", token.INT, uint64(12)},
		{`""`, token.STRING, ""},
		{`"1"`, token.STRING, "1"},

		{`"12"`, token.STRING, "12"}, // 15
		{"``", token.STRING, ""},
		{"`1`", token.STRING, "1"},
		{"`12`", token.STRING, "12"},
		{"'@'", token.CHAR, int32('@')},

		{"a ", token.IDENT, "a"}, // 20
		{"ab ", token.IDENT, "ab"},
		{"1 ", token.INT, uint64(1)},
		{"12 ", token.INT, uint64(12)},
		{`"" `, token.STRING, ""},

		{`"1" `, token.STRING, "1"}, // 25
		{`"12" `, token.STRING, "12"},
		{"`` ", token.STRING, ""},
		{"`1` ", token.STRING, "1"},
		{"`12` ", token.STRING, "12"},

		{"'@' ", token.CHAR, int32('@')}, // 30
		{" a", token.IDENT, "a"},
		{" ab", token.IDENT, "ab"},
		{" 1", token.INT, uint64(1)},
		{" 12", token.INT, uint64(12)},

		{` ""`, token.STRING, ""}, // 35
		{` "1"`, token.STRING, "1"},
		{` "12"`, token.STRING, "12"},
		{" ``", token.STRING, ""},
		{" `1`", token.STRING, "1"},

		{" `12`", token.STRING, "12"}, // 40
		{" '@'", token.CHAR, int32('@')},
		{" a ", token.IDENT, "a"},
		{" ab ", token.IDENT, "ab"},
		{" 1 ", token.INT, uint64(1)},

		{" 12 ", token.INT, uint64(12)}, // 45
		{` "" `, token.STRING, ""},
		{` "1" `, token.STRING, "1"},
		{` "12" `, token.STRING, "12"},
		{" `` ", token.STRING, ""},

		{" `1` ", token.STRING, "1"}, // 50
		{" `12` ", token.STRING, "12"},
		{" '@' ", token.CHAR, int32('@')},
		{"bá", token.IDENT, "bá"},
		{"bár", token.IDENT, "bár"},

		{"bára", token.IDENT, "bára"}, // 55
	})
}

func test(t *testing.T, root string) {
	var (
		count    int
		tokCount int
		size     int64
	)

	if err := filepath.Walk(
		root,
		func(pth string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			ok, err := filepath.Match("*.go", filepath.Base(pth))
			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				return nil
			}

			q := pth[len(root)+1:]
			for _, v := range whitelist {
				if q == v {
					return nil
				}
			}

			file, err := os.Open(pth)
			if err != nil {
				t.Fatal(err)
			}

			defer file.Close()
			src, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}

			var sc goscanner.Scanner
			fs := token.NewFileSet()
			base := fs.Base()
			fl := fs.AddFile(pth, base, int(info.Size()))
			sc.Init(fl, src, nil, goscanner.ScanComments)

			l := New(src)

			i := 1
			//dbg("pth %q", pth)
		loop:
			for {
				pos, tok, lit := sc.Scan()
				if tok == token.SEMICOLON && lit == "\n" { // injected
					continue
				}

				tok2, lit2 := l.Scan()
				pos2 := token.Pos(base + l.i0 - 1)
				p, p2 := fs.Position(pos), fs.Position(pos2)

				if g, e := tok2, tok; g != e {
					t.Errorf(
						"%d.%d got tok %s(%d), exp tok %s(%d) got lit '%v', exp lit '%v' %s %s",
						count, i, g, int(g), e, int(e), lit2, lit, p2, p,
					)

					if tok2 == token.ILLEGAL && len(l.Errors) != 0 {
						t.Fatal(l.Errors[0], p2, p)
					}
					return nil
				}

				if g, e := p2, p; g != e {
					t.Fatal(count, i, tok, g, e)
				}

				if lit2 == nil {
					lit2 = tok2.String()
				}
				switch tok2 {
				case token.EOF:
					break loop
				case token.CHAR:
					s, _ := strconv.Unquote(lit)
					if s == "" {
						s = lit
					}
					switch x := lit2.(type) {
					case int32:
						if g, e := x, int32(s[0]); g != e {
							t.Fatalf("%d.%d %s %v %v %s %s", count, i, tok, g, e, p2, p)
						}
					case string:
						if g, e := x, lit; g != e {
							t.Errorf("%d.%d %s g: %q e: %q %s %s", count, i, tok, g, e, p2, p)
						}
					default:
						t.Fatalf("%d: %T(%#v)", i, lit2, lit2)
					}
				case token.INT:
					n, _ := strconv.ParseUint(lit, 0, 64)
					switch x := lit2.(type) {
					case uint64:
						if g, e := x, n; g != e {
							t.Fatalf("%d.%d %s %v %v %s %s", count, i, tok, g, e, p2, p)
						}
					case string:
						if g, e := x, lit; g != e {
							t.Fatalf("%d.%d %s %q %q %s %s", count, i, tok, g, e, p2, p)
						}
					default:
						t.Fatalf("%d: %T(%#v)", i, lit2, lit2)
					}
				case token.IMAG:
					n, _ := strconv.ParseFloat(lit[:len(lit)-1], 64)
					switch x := lit2.(type) {
					case complex128:
						if g, e := x, complex(0, n); g != e {
							t.Fatalf("%d.%d %s %v %v %s %s", count, i, tok, g, e, p2, p)
						}
					case string:
						if g, e := x, lit; g != e {
							t.Fatalf("%d.%d %s %q %q %s %s", count, i, tok, g, e, p2, p)
						}
					default:
						t.Fatalf("%d: %T(%#v)", i, lit2, lit2)
					}
				case token.FLOAT:
					n, _ := strconv.ParseFloat(lit, 64)
					switch x := lit2.(type) {
					case float64:
						if g, e := x, n; g != e {
							t.Fatalf("%d.%d %s %v %v %s %s", count, i, tok, g, e, p2, p)
						}
					case string:
						if g, e := x, lit; g != e {
							t.Fatalf("%d.%d %s %q %q %s %s", count, i, tok, g, e, p2, p)
						}
					default:
						t.Fatalf("%d: %T(%#v)", i, lit2, lit2)
					}
				case token.STRING:
					nlit, err := strconv.Unquote(lit)
					if err != nil {
						nlit = lit
					}
					if g, e := lit2.(string), nlit; g != e {
						t.Fatalf("%d.%d %s %q %q %s %s // %q", count, i, tok, g, e, p2, p, lit)
					}
				case token.ILLEGAL:
					if g, e := lit2.(string), lit; g != e {
						t.Fatalf("%d.%d %s %q %q %s %s", count, i, tok, g, e, p2, p)
					}
				}

				tokCount++
				i++
			}

			count++
			size += info.Size()
			return nil
		}); err != nil {
		t.Fatal(err)
	}

	t.Logf("%d .go files, %d bytes, %d tokens\n", count, size, tokCount)
}

func TestStdlib(t *testing.T) {
	test(t, std)
}

func TestTestData(t *testing.T) {
	test(t, "testdata")
}

func TestTests(t *testing.T) {
	test(t, tests)
}

func testSemis(t *testing.T, root string) {
	var (
		count    int
		tokCount int
		size     int64
	)

	if err := filepath.Walk(
		root,
		func(pth string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			ok, err := filepath.Match("*.go", filepath.Base(pth))
			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				return nil
			}

			q := pth[len(root)+1:]
			for _, v := range whitelist {
				if q == v {
					return nil
				}
			}

			file, err := os.Open(pth)
			if err != nil {
				t.Fatal(err)
			}

			defer file.Close()
			src, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}

			var sc goscanner.Scanner
			fs := token.NewFileSet()
			base := fs.Base()
			fl := fs.AddFile(pth, base, int(info.Size()))
			sc.Init(fl, src, nil, goscanner.ScanComments)

			l := New(src)

			i := 1
			//dbg("pth %q", pth)
		loop:
			for {
				pos, tok, lit := sc.Scan()
				tok2, lit2 := l.ScanSemis()
				pos2 := token.Pos(base + l.i0 - 1)
				p, p2 := fs.Position(pos), fs.Position(pos2)

				if g, e := tok2, tok; g != e {
					t.Fatalf(
						"%d.%d got tok %s(%d), exp tok %s(%d) got lit '%v', exp lit '%v' %s %s",
						count, i, g, int(g), e, int(e), lit2, lit, p2, p,
					)

					//if tok2 == token.ILLEGAL && len(l.Errors) != 0 {
					//	t.Fatal(l.Errors[0], p2, p)
					//}
					//return nil
				}

				if g, e := p2, p; g != e {
					t.Fatal(count, i, tok, g, e)
				}

				if lit2 == nil {
					lit2 = tok2.String()
				}
				switch tok2 {
				case token.EOF:
					break loop
				case token.CHAR:
					s, _ := strconv.Unquote(lit)
					if s == "" {
						s = lit
					}
					switch x := lit2.(type) {
					case int32:
						if g, e := x, int32(s[0]); g != e {
							t.Fatalf("%d.%d %s %v %v %s %s", count, i, tok, g, e, p2, p)
						}
					case string:
						if g, e := x, lit; g != e {
							t.Errorf("%d.%d %s %q %q %s %s", count, i, tok, g, e, p2, p)
						}
					default:
						t.Fatalf("%d: %T(%#v)", i, lit2, lit2)
					}
				case token.INT:
					n, _ := strconv.ParseUint(lit, 0, 64)
					switch x := lit2.(type) {
					case uint64:
						if g, e := x, n; g != e {
							t.Fatalf("%d.%d %s %v %v %s %s", count, i, tok, g, e, p2, p)
						}
					case string:
						if g, e := x, lit; g != e {
							t.Fatalf("%d.%d %s %q %q %s %s", count, i, tok, g, e, p2, p)
						}
					default:
						t.Fatalf("%d: %T(%#v)", i, lit2, lit2)
					}
				case token.IMAG:
					n, _ := strconv.ParseFloat(lit[:len(lit)-1], 64)
					switch x := lit2.(type) {
					case complex128:
						if g, e := x, complex(0, n); g != e {
							t.Fatalf("%d.%d %s %v %v %s %s", count, i, tok, g, e, p2, p)
						}
					case string:
						if g, e := x, lit; g != e {
							t.Fatalf("%d.%d %s %q %q %s %s", count, i, tok, g, e, p2, p)
						}
					default:
						t.Fatalf("%d: %T(%#v)", i, lit2, lit2)
					}
				case token.FLOAT:
					n, _ := strconv.ParseFloat(lit, 64)
					switch x := lit2.(type) {
					case float64:
						if g, e := x, n; g != e {
							t.Fatalf("%d.%d %s %v %v %s %s", count, i, tok, g, e, p2, p)
						}
					case string:
						if g, e := x, lit; g != e {
							t.Fatalf("%d.%d %s %q %q %s %s", count, i, tok, g, e, p2, p)
						}
					default:
						t.Fatalf("%d: %T(%#v)", i, lit2, lit2)
					}
				case token.STRING:
					nlit, err := strconv.Unquote(lit)
					if err != nil {
						nlit = lit
					}
					if g, e := lit2.(string), nlit; g != e {
						t.Fatalf("%d.%d %s %q %q %s %s // %q", count, i, tok, g, e, p2, p, lit)
					}
				case token.ILLEGAL:
					if g, e := lit2.(string), lit; g != e {
						t.Fatalf("%d.%d %s %q %q %s %s", count, i, tok, g, e, p2, p)
					}
				}

				tokCount++
				i++
			}

			count++
			size += info.Size()
			return nil
		}); err != nil {
		t.Fatal(err)
	}

	t.Logf("%d .go files, %d bytes, %d tokens\n", count, size, tokCount)
}

func TestStdlibSemis(t *testing.T) {
	testSemis(t, std)
}

func TestTestDataSemis(t *testing.T) {
	testSemis(t, "testdata")
}

func TestTestsSemis(t *testing.T) {
	testSemis(t, tests)
}

func TestBenchGoScanner(t *testing.T) {
	t0 := time.Now()
	if err := filepath.Walk(
		std,
		func(pth string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			ok, err := filepath.Match("*.go", filepath.Base(pth))
			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				return nil
			}

			file, err := os.Open(pth)
			if err != nil {
				t.Fatal(err)
			}

			defer file.Close()
			src, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}

			var sc goscanner.Scanner
			fs := token.NewFileSet()
			base := fs.Base()
			fl := fs.AddFile(pth, base, int(info.Size()))
			sc.Init(fl, src, nil, goscanner.ScanComments)
			for {
				_, tok, _ := sc.Scan()
				if tok == token.EOF {
					break
				}
			}

			return nil
		}); err != nil {
		t.Fatal(err)
	}

	t.Log(time.Since(t0))
}

func TestBenchScanner(t *testing.T) {
	t0 := time.Now()
	if err := filepath.Walk(
		std,
		func(pth string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			ok, err := filepath.Match("*.go", filepath.Base(pth))
			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				return nil
			}

			file, err := os.Open(pth)
			if err != nil {
				t.Fatal(err)
			}

			defer file.Close()
			src, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}

			l := New(src)
			for {
				tok, _ := l.Scan()
				if tok == token.EOF {
					break
				}
			}

			return nil
		}); err != nil {
		t.Fatal(err)
	}

	t.Log(time.Since(t0))
}
