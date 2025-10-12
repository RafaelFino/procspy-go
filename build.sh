#!/bin/bash
set -euo pipefail

# if no args, consider "all"
par=${1:-all}

if [ "$par" == "clean" ]; then
    echo "Cleaning bin directory"
    rm -rf bin
    exit 0
fi

if [ "$par" == "all" ]; then
    archs=( "amd64" )
    oses=( "linux" "windows" )

    for os in "${oses[@]}"; do
        for arch in "${archs[@]}"; do
            for d in cmd/* ; do
                name=${d##*/}
                outdir=bin/$os-$arch
                mkdir -p "$outdir"
                echo "[$os $arch] Building $name -> $outdir/$name"
                GOOS=$os GOARCH=$arch CGO_ENABLED=1 \
                    go build -ldflags "-s -w -X main.buildDate=$(date -u +'%Y-%m-%d_%H:%M:%S')" \
                    -o "$outdir/$name" "$d/main.go"
            done
        done
    done
    exit 0
fi

os=$(go env GOOS)
arch=$(go env GOARCH)

for d in cmd/* ; do
    name=${d##*/}
    outdir=bin/$os-$arch
    mkdir -p "$outdir"
    echo "[$os $arch] Building $name -> $outdir/$name"
    GOOS=$os GOARCH=$arch CGO_ENABLED=1 \
        go build -o "$outdir/$name" "$d/main.go"
done
