#!/usr/bin/env sh

for os in $(echo darwin linux); do
  for arch in $(echo amd64 386); do
    echo "==> Creating ${os}_${arch} build."

    GOOS=$os GOARCH=$arch go build -o dist/gadget_${os}_${arch} .
  done
done
