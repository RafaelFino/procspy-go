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
            echo " >> Building CLI for $os $arch > bin/$os-$arch/procspy"
            GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-cli cmd/cli/main.go

            echo " >> Building Client for $os $arch> bin/$os-$arch/procspy-client"
            GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-client cmd/client/main.go

            echo " >> Building Server for $os $arch > bin/$os-$arch/procspy-server"
            GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-server cmd/server/main.go
        done
    done
    exit 0    
fi

os=`go env GOOS`
arch=`go env GOARCH`

echo " >> Building CLI for $os $arch > bin/$os-$arch/procspy"
GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-cli cmd/cli/main.go

echo " >> Building Client for $os $arch> bin/$os-$arch/procspy-client"
GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-client cmd/client/main.go

echo " >> Building Server for $os $arch > bin/$os-$arch/procspy-server"
GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-server cmd/server/main.go