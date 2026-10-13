// Command passphrase generates random diceware-style passphrases using the
// EFF long wordlist and a cryptographically secure random source.
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("passphrase", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: passphrase [flags]\n\n")
		fmt.Fprintf(stderr, "Generate random passphrases from the EFF long wordlist (7776 words).\n\n")
		fmt.Fprintf(stderr, "Flags:\n")
		fs.PrintDefaults()
	}

	var opts Options
	var count int
	var showEntropy bool
	fs.IntVar(&opts.Words, "words", 6, "number of words per passphrase")
	fs.StringVar(&opts.Separator, "sep", "-", "separator placed between words")
	fs.BoolVar(&opts.Capitalize, "capitalize", false, "capitalize the first letter of each word")
	fs.IntVar(&count, "count", 1, "number of passphrases to generate")
	fs.BoolVar(&showEntropy, "entropy", false, "print the entropy of each passphrase to stderr")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "passphrase: unexpected argument %q\n", fs.Arg(0))
		fs.Usage()
		return 2
	}
	if count < 1 {
		fmt.Fprintln(stderr, "passphrase: -count must be at least 1")
		return 2
	}
	if opts.Words < 1 {
		fmt.Fprintln(stderr, "passphrase: -words must be at least 1")
		return 2
	}

	for range count {
		p, err := Generate(rand.Reader, opts)
		if err != nil {
			fmt.Fprintf(stderr, "passphrase: %v\n", err)
			return 1
		}
		fmt.Fprintln(stdout, p)
	}
	if showEntropy {
		fmt.Fprintf(stderr, "%.1f bits of entropy per passphrase (%d words from a %d-word list)\n",
			Entropy(opts.Words), opts.Words, len(wordlist))
	}
	return 0
}
