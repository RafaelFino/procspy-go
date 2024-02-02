#!/bin/bash
par=$1

if [ "$par" == "clean" ]; then
    rm -rf bin
    exit 0
fi

if [ "$par" == "all" ]; then
    archs=( "amd64" )
    oses=( "linux" "windows" )

    for os in "${oses[@]}"
    do
        for arch in "${archs[@]}"
        do
            echo "Building for $os $arch"
            GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy cmd/procspy.go
        done
    done
    exit 0    
fi

os=`go env GOOS`
arch=`go env GOARCH`
go build -o bin/$os-$arch/procspy cmd/procspy.go

