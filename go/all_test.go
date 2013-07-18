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

func dbg(s string, va ...interface{}) {
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Printf("%s:%d: ", path.Base(fn), fl)
	fmt.Printf(s, va...)
	fmt.Println()
}

var std = filepath.Join(runtime.GOROOT(), "/src")

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
			t.Fatal(i, g, e)
		}

		if g, e := lit, test.lit; g != e {
			t.Fatal(i, g, e)
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
					if tok2 == token.ILLEGAL && len(l.Errors) != 0 {
						t.Fatal(l.Errors[0], p2, p)
					}

					t.Fatalf("%d.%d %s(%d) %s(%d) '%v' '%v' %s %s", count, i, g, int(g), e, int(e), lit2, lit, p2, p)

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
					if g, e := lit2.(int32), int32(s[0]); g != e {
						t.Fatalf("%d.%d %s %v %v %s %s", count, i, tok, g, e, p2, p)
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
					}
				case token.FLOAT:
					n, _ := strconv.ParseFloat(lit, 64)
					if g, e := lit2.(float64), n; g != e {
						t.Fatalf("%d.%d %s %v %v %s %s", count, i, tok, g, e, p2, p)
					}
				case token.STRING:
					lit, _ = strconv.Unquote(lit)
					if g, e := lit2.(string), lit; g != e {
						t.Fatalf("%d.%d %s %q %q %s %s", count, i, tok, g, e, p2, p)
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

func Test0(t *testing.T) {
	test(t, std)
}

func Test1(t *testing.T) {
	test(t, "testdata")
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
