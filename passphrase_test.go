package main

import (
	"bytes"
	"crypto/rand"
	"strings"
	"testing"
)

func TestWordlist(t *testing.T) {
	if got := len(wordlist); got != 7776 {
		t.Fatalf("len(wordlist) = %d, want 7776", got)
	}
	seen := make(map[string]bool, len(wordlist))
	for _, w := range wordlist {
		if w == "" {
			t.Fatal("wordlist contains an empty word")
		}
		if seen[w] {
			t.Fatalf("wordlist contains duplicate word %q", w)
		}
		seen[w] = true
	}
	if wordlist[0] != "abacus" || wordlist[len(wordlist)-1] != "zoom" {
		t.Errorf("wordlist runs %q..%q, want abacus..zoom", wordlist[0], wordlist[len(wordlist)-1])
	}
}

func TestGenerate(t *testing.T) {
	inList := make(map[string]bool, len(wordlist))
	for _, w := range wordlist {
		inList[w] = true
	}

	p, err := Generate(rand.Reader, Options{Words: 5, Separator: " "})
	if err != nil {
		t.Fatal(err)
	}
	words := strings.Split(p, " ")
	if len(words) != 5 {
		t.Fatalf("Generate returned %d words, want 5: %q", len(words), p)
	}
	for _, w := range words {
		if !inList[w] {
			t.Errorf("word %q is not in the wordlist", w)
		}
	}
}

func TestGenerateCapitalize(t *testing.T) {
	p, err := Generate(rand.Reader, Options{Words: 4, Separator: ".", Capitalize: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range strings.Split(p, ".") {
		if w[:1] != strings.ToUpper(w[:1]) {
			t.Errorf("word %q is not capitalized", w)
		}
	}
}

func TestGenerateRejectsZeroWords(t *testing.T) {
	if _, err := Generate(rand.Reader, Options{Words: 0}); err == nil {
		t.Fatal("Generate with 0 words succeeded, want error")
	}
}

func TestEntropy(t *testing.T) {
	// log2(7776) is about 12.925 bits per word.
	if got := Entropy(6); got < 77.5 || got > 77.6 {
		t.Errorf("Entropy(6) = %v, want ~77.55", got)
	}
}

func TestRun(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"-count", "3", "-words", "4", "-sep", "_"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run exited %d, stderr: %s", code, stderr.String())
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("got %d passphrases, want 3: %q", len(lines), stdout.String())
	}
	for _, l := range lines {
		if n := len(strings.Split(l, "_")); n != 4 {
			t.Errorf("passphrase %q has %d words, want 4", l, n)
		}
	}
}

func TestRunInvalidArgs(t *testing.T) {
	for _, args := range [][]string{
		{"-words", "0"},
		{"-count", "0"},
		{"extra"},
		{"-nope"},
	} {
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr); code != 2 {
			t.Errorf("run(%q) exited %d, want 2", args, code)
		}
	}
}
