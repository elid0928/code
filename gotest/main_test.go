package gotest

import "testing"

func FuzzTest(f *testing.F) {
	cases := []string{
		"test",
		"tset",
		"foo",
		"bar",
	}

	for _, c := range cases {
		f.Add(c)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		reversed := Reverse(orig)
		if reversed == orig {
			t.Errorf("Reverse(%q) = %q, want anything but %q", orig, reversed, orig)
		}
	})
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
