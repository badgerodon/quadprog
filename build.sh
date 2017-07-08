#!/bin/sh

OS=$(uname -s | tr '[:upper:]' '[:lower:]')

ls src -1 \
    | xargs -I {} basename {} .f \
    | xargs -I {} gfortran -fPIC -c -o /tmp/{}.o src/{}.f

ls src -1 \
    | xargs -I {} basename {} .f \
    | xargs -I {} echo /tmp/{}.o \
    | xargs ld -r -o /tmp/quadprog_${OS}_amd64.syso

ls -alht /tmp

if [ "$OS" == "darwin" ]; then
    docker build . -f build-linux.dockerfile -t badgerodon-quadprog-build
    docker create --name extract badgerodon-quadprog-build:latest
    docker cp extract:/tmp/quadprog_linux_amd64.syso quadprog_linux_amd64.syso
    docker rm -f extract
fi
