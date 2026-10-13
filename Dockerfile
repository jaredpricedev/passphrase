# syntax=docker/dockerfile:1

# Build stage: cross-compile a static binary for the target platform.
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY go.mod ./
COPY *.go eff_large_wordlist.txt ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/passphrase .

# Runtime stage: nothing but the binary, running as an unprivileged user.
FROM scratch
COPY --from=build /out/passphrase /passphrase
USER 65534:65534
ENTRYPOINT ["/passphrase"]
