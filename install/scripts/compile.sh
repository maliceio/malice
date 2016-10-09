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
rm -rf build/*
rm -rf release/*
mkdir -p release

# Get the git commit
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$(git describe --tags)

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"arm amd64"}
XC_OS=${XC_OS:-"solaris darwin freebsd linux"}

# Build!
echo "==> Building linux..."
for GOARCH in arm amd64; do
go build -v -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.GitDescribe=${GIT_DESCRIBE}" -o build/linux_${GOARCH}/malice
done

echo "==> Building darwin..."
go build -v -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.GitDescribe=${GIT_DESCRIBE}" -o build/darwin_amd64/malice

# Zip and copy to the dist dir
echo "==> Packaging..."
for PLATFORM in $(find ./build -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "--> ${OSARCH}"

    pushd $PLATFORM >/dev/null 2>&1
    zip ../${OSARCH}.zip ./*
    popd >/dev/null 2>&1
done

# Done!
echo
echo "==> Results:"
ls -hl build/