#!/bin/bash
#
# This script builds the application from source for multiple platforms.
set -e

export CGO_ENABLED=0

# Install dependencies
echo "==> Getting dependencies..."
go get github.com/mitchellh/gox
go get -v -d ./...

# Delete the old dir
echo "==> Removing old directory..."
rm -rf build/*
rm -rf release/*
mkdir -p release

# Get the git commit
GIT_COMMIT="$(git rev-parse --short HEAD)"
GIT_DIRTY="$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"
GIT_DESCRIBE="$(git describe --tags --always)"

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
XC_OS=${XC_OS:-"solaris darwin freebsd linux"}

# Build!
echo "==> Building..."
"`which gox`" \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -osarch="!darwin/arm !solaris/amd64" \
    -ldflags "-X main.GitCommit='${GIT_COMMIT}${GIT_DIRTY}' -X main.GitDescribe='${GIT_DESCRIBE}'" \
    -output "build/{{.OS}}_{{.Arch}}/malice" \
    .

# Zip and copy to the dist dir
echo "==> Packaging..."
for PLATFORM in $(find ./build -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "--> ${OSARCH}"

    pushd $PLATFORM >/dev/null 2>&1
    zip ../${NAME}_${VERSION}_${OSARCH}.zip ./*
    popd >/dev/null 2>&1
done

# Done!
echo
echo "==> Results:"
ls -hl build/
