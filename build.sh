#!/bin/bash
archs=( "arm64" "amd64" )
oses=( "linux" "windows" "darwin" )

for os in "${oses[@]}"
do
    for arch in "${archs[@]}"
    do
        echo "Building for $os $arch"
        GOOS=$os GOARCH=$arch go build -o bin/$os-$arch/procspy procspy.go
    done
done