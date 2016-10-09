#!/bin/bash
#
# This script builds the application from source for multiple platforms.
set -e

go get github.com/mitchellh/gox

# Install dependencies
echo "==> Getting dependencies..."
# go get ./...
go get -v -d ./...

# Delete the old dir
echo "==> Removing old directory..."
rm -f bin/*
rm -rf pkg/*
mkdir -p bin/

# Get the git commit
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$(git describe --tags)

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"386 amd64"}
XC_OS=${XC_OS:-"solaris darwin freebsd linux"}

# Build!
echo "==> Building..."
for GOOS in darwin linux; do
  for GOARCH in 386 amd64; do
    go build -v -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.GitDescribe=${GIT_DESCRIBE}" -o build/$GOOS_$GOARCH/malice
  done
done
# gox \
#     -os="${XC_OS}" \
#     -arch="${XC_ARCH}" \
#     -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.GitDescribe=${GIT_DESCRIBE}" \
#     -output "build/{{.OS}}_{{.Arch}}/malice" \
#     .

