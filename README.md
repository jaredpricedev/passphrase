# passphrase

A tiny command-line passphrase generator, packaged as a Docker image.

It picks words uniformly at random from the [EFF long wordlist](https://www.eff.org/deeplinks/2016/07/new-wordlists-random-passphrases)
(7776 words, about 12.9 bits of entropy per word) using Go's `crypto/rand`.
The default of 6 words gives about 77.5 bits of entropy.

```
$ docker run --rm ghcr.io/jaredpricedev/passphrase
glandular-swimming-rebound-grumble-flyaway-clatter
```

## Run it

### Use the prebuilt image

GitHub Actions builds the image for `linux/amd64` and `linux/arm64` (so it also
runs on Apple Silicon Macs) and publishes it to GitHub Container Registry:

```sh
docker run --rm ghcr.io/jaredpricedev/passphrase
```

Tags: `latest` (the `main` branch), one tag per branch (`/` replaced by `-`),
`sha-<commit>`, and `X.Y.Z` for `vX.Y.Z` git tags.

### Build it yourself

```sh
git clone https://github.com/jaredpricedev/passphrase.git
cd passphrase
docker build -t passphrase .
docker run --rm passphrase
```

The final image is built `FROM scratch` and contains only the static binary
(about 3 MB). It runs as an unprivileged user and needs no network access, so
you can add `--network none` if you like.

## Options

```
-words int      number of words per passphrase (default 6)
-count int      number of passphrases to generate (default 1)
-sep string     separator placed between words (default "-")
-capitalize     capitalize the first letter of each word
-entropy        print the entropy of each passphrase to stderr
```

Examples:

```sh
docker run --rm passphrase -count 5
docker run --rm passphrase -words 8 -sep ' ' -capitalize
docker run --rm passphrase -words 4 -sep . -entropy
```

Tip: add a shell alias so it feels like a local command:

```sh
alias passphrase='docker run --rm --network none ghcr.io/jaredpricedev/passphrase'
```

## Develop

```sh
go test ./...
go run . -count 3
```

## Credits

The wordlist (`eff_large_wordlist.txt`) is by the
[Electronic Frontier Foundation](https://www.eff.org/dice), licensed under
[CC BY 3.0 US](https://creativecommons.org/licenses/by/3.0/us/).
