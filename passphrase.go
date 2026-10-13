package main

import (
	"bufio"
	"crypto/rand"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"strings"
	"unicode"
	"unicode/utf8"
)

// effLargeWordlist is the EFF long wordlist (7776 words, one per line in the
// form "<dice roll>\t<word>").
// Source: https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt
// License: CC BY 3.0 US, Electronic Frontier Foundation.
//
//go:embed eff_large_wordlist.txt
var effLargeWordlist string

// wordlist holds the words parsed from effLargeWordlist.
var wordlist = mustParseWordlist(effLargeWordlist)

func mustParseWordlist(raw string) []string {
	var words []string
	sc := bufio.NewScanner(strings.NewReader(raw))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		// Lines look like "11111\tabacus"; keep only the word.
		fields := strings.Fields(line)
		words = append(words, fields[len(fields)-1])
	}
	if len(words) != 7776 {
		panic(fmt.Sprintf("wordlist: expected 7776 words, got %d", len(words)))
	}
	return words
}

// Options controls how a passphrase is generated.
type Options struct {
	Words      int
	Separator  string
	Capitalize bool
}

// Generate returns a passphrase made of opts.Words words chosen uniformly at
// random from the wordlist, reading randomness from rnd (normally
// crypto/rand.Reader).
func Generate(rnd io.Reader, opts Options) (string, error) {
	if opts.Words < 1 {
		return "", errors.New("number of words must be at least 1")
	}
	max := big.NewInt(int64(len(wordlist)))
	parts := make([]string, opts.Words)
	for i := range parts {
		n, err := rand.Int(rnd, max)
		if err != nil {
			return "", fmt.Errorf("reading random data: %w", err)
		}
		word := wordlist[n.Int64()]
		if opts.Capitalize {
			word = capitalize(word)
		}
		parts[i] = word
	}
	return strings.Join(parts, opts.Separator), nil
}

// Entropy returns the entropy in bits of a passphrase with the given number
// of words.
func Entropy(words int) float64 {
	return float64(words) * math.Log2(float64(len(wordlist)))
}

func capitalize(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}
