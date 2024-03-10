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
            echo "[$os $arch] Building CLI -> ./bin/$os-$arch/procspy"
            GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-cli cmd/cli/main.go

            echo "[$os $arch] Building Client -> ./bin/$os-$arch/procspy-client"
            GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-client cmd/client/main.go

            echo "[$os $arch] Building Server -> ./bin/$os-$arch/procspy-server"
            GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-server cmd/server/main.go
        done
    done
    exit 0    
fi

os=`go env GOOS`
arch=`go env GOARCH`

echo "[$os $arch] Building CLI -> ./bin/$os-$arch/procspy"
GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-cli cmd/cli/main.go

echo "[$os $arch] Building Client -> ./bin/$os-$arch/procspy-client"
GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-client cmd/client/main.go

echo "[$os $arch] Building Server -> ./bin/$os-$arch/procspy-server"
GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -o bin/$os-$arch/procspy-server cmd/server/main.go