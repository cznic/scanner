// Copyright (c) 2014 Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"testing"
	"unicode"
)

func caller(s string, va ...interface{}) {
	_, fn, fl, _ := runtime.Caller(2)
	fmt.Printf("caller: %s:%d: ", path.Base(fn), fl)
	fmt.Printf(s, va...)
	fmt.Println()
	_, fn, fl, _ = runtime.Caller(1)
	fmt.Printf("\tcallee: %s:%d: ", path.Base(fn), fl)
	fmt.Println()
}

func dbg(s string, va ...interface{}) {
	if s == "" {
		s = strings.Repeat("%v ", len(va))
	}
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Printf("dbg %s:%d: ", path.Base(fn), fl)
	fmt.Printf(s, va...)
	fmt.Println()
}

func TODO(...interface{}) string {
	_, fn, fl, _ := runtime.Caller(1)
	return fmt.Sprintf("TODO: %s:%d:\n", path.Base(fn), fl)
}

func use(...interface{}) {}

// ============================================================================

func TestScanner(t *testing.T) {
	tab := []struct {
		src string
		ok  bool
		tok Token
		val string
	}{
		// 0
		{"_", true, ILLEGAL, "_"},
		{"_ ", true, ILLEGAL, "_"},
		{"_\n", true, ILLEGAL, "_"},
		{"__", true, ILLEGAL, "_"},
		{"__ ", true, ILLEGAL, "_"},

		// 5
		{"__\n", true, ILLEGAL, "_"},
		{"_:", true, ILLEGAL, "_"},
		{"_: ", true, ILLEGAL, "_"},
		{"_:\n", true, ILLEGAL, "_"},
		{"_:!", true, ILLEGAL, "_"},

		// 10
		{"_:0", true, LABEL, "0"},
		{" _:0", true, LABEL, "0"},
		{"\t_:0\t", true, LABEL, "0"},
		{"\n_:0\n", true, LABEL, "0"},
		{"\n\t_:0\t\n", true, LABEL, "0"},
	}

	for i, test := range tab {
		sc := New([]byte(test.src))
		tok, val := sc.Scan()
		errs := sc.Errors
		switch test.ok {
		case true:
			if len(errs) != 0 {
				t.Error("errs", i, errs)
				break
			}

			if g, e := tok, test.tok; g != e {
				t.Error("tok", i, g, e)
			}

			if g, e := val, test.val; g != e {
				t.Errorf("val %v %q %q", i, g, e)
			}
		default:
			if len(errs) == 0 {
				t.Error("errs", i, tok, val)
				break
			}
		}
	}
}

func encodeRune(r rune) string {
	switch i := uint32(r); {
	case i <= rune1Max:
		return string([]byte{byte(r)})
	case i <= rune2Max:
		return string([]byte{t2 | byte(r>>6), tx | byte(r)&maskx})
	case i <= rune3Max:
		return string([]byte{t3 | byte(r>>12), tx | byte(r>>6)&maskx, tx | byte(r)&maskx})
	default:
		return string([]byte{t4 | byte(r>>18), tx | byte(r>>12)&maskx, tx | byte(r>>6)&maskx, tx | byte(r)&maskx})
	}
}

func TestLabel(t *testing.T) {
	for c := rune(0); c <= unicode.MaxRune; c++ {
		s := "_:" + encodeRune(c)
		sc := New([]byte(s))
		tok, val := sc.Scan()
		switch {
		case c >= '0' && c <= '9', checkPnCharsU(c):
			if g, e := tok, LABEL; g != e {
				t.Fatalf("tok: %q(%U) %v %v", c, c, g, e)
			}
			if g, e := val, s[2:]; g != e {
				t.Fatalf("val: %q(%U) %v %v", c, c, g, e)
			}
		default:
			if g, e := tok, ILLEGAL; g != e {
				t.Fatalf("should fail %q |% x|: %q(%U) %v %v", s, s, c, c, g, e)
			}
		}
	}

	for c := rune(1); c <= unicode.MaxRune; c++ {
		s := "_:0" + encodeRune(c)
		sc := New([]byte(s))
		tok, val := sc.Scan()
		switch {
		case checkPnChars(c):
			if g, e := tok, LABEL; g != e {
				t.Fatalf("tok: %q(%U) %v %v", c, c, g, e)
			}
			if g, e := val, s[2:]; g != e {
				t.Fatalf("val: %q(%U) %v %v", c, c, g, e)
			}
		case c == '.', c >= 0x80:
			if g, e := tok, ILLEGAL; g != e {
				t.Fatalf("should fail %q |% x|: %q(%U) %v %v", s, s, c, c, g, e)
			}
		default:
			if g, e := tok, LABEL; g != e {
				t.Fatalf("tok: %q(%U) %v %v", c, c, g, e)
			}
			if g, e := val, "0"; g != e {
				t.Fatalf("val: %q(%U) %q %q", c, c, g, e)
			}
		}
	}

	for c := rune(1); c <= unicode.MaxRune; c++ {
		s := "_:0a" + encodeRune(c)
		sc := New([]byte(s))
		tok, val := sc.Scan()
		switch {
		case checkPnChars(c):
			if g, e := tok, LABEL; g != e {
				t.Fatalf("tok: %q(%U) %v %v", c, c, g, e)
			}
			if g, e := val, s[2:]; g != e {
				t.Fatalf("val: %q(%U) %v %v", c, c, g, e)
			}
		case c == '.', c >= 0x80:
			if g, e := tok, ILLEGAL; g != e {
				t.Fatalf("should fail %q |% x|: %q(%U) %v %v", s, s, c, c, g, e)
			}
		default:
			if g, e := tok, LABEL; g != e {
				t.Fatalf("tok: %q(%U) %v %v", c, c, g, e)
			}
			if g, e := val, "0a"; g != e {
				t.Fatalf("val: %q(%U) %q %q", c, c, g, e)
			}
		}
	}

	for c := rune(1); c <= unicode.MaxRune; c++ {
		s := "_:0." + encodeRune(c)
		sc := New([]byte(s))
		tok, val := sc.Scan()
		switch {
		case checkPnChars(c):
			if g, e := tok, LABEL; g != e {
				t.Fatalf("tok: %q(%U) %v %v", c, c, g, e)
			}
			if g, e := val, s[2:]; g != e {
				t.Fatalf("val: %q(%U) %v %v", c, c, g, e)
			}
		default:
			if g, e := tok, ILLEGAL; g != e {
				t.Fatalf("should fail %q |% x|: %q(%U) %v %v", s, s, c, c, g, e)
			}
		}
	}
}
