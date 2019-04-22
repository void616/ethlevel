#!/bin/bash

case "$(uname -s)" in
    Linux*)     machine=linux;;
    Darwin*)    machine=mac;;
    CYGWIN*)    machine=cygwin;;
    MINGW*)     machine=mingw;;
    *)          machine="unknown"
esac

if [ "$machine" == "cygwin" ] || [ "$machine" == "mingw" ]; then
	MOUNT_GOPATH="/${GOPATH//\\/$'/'}"
	MOUNT_PWD="/${PWD//\\/$'/'}"
else
	MOUNT_GOPATH=$GOPATH
	MOUNT_PWD=$(PWD)
fi

docker run \
	--rm \
	-i --name gobuild_with_docker \
	-e "GOOS=$BUILD_OS" \
	-e "GOARCH=$BUILD_ARCH" \
	-e "CGO=1" \
	-v "$MOUNT_GOPATH":/go \
	-v "$MOUNT_PWD":/app \
	gobuild_with_docker \
	go build -ldflags "-s" -tags "$BUILD_TAGS" -o "$BUILD_OUTDIR/$BUILD_OUTFILE" "$BUILD_APP"