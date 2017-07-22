#!/bin/sh

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [ "$(uname -m)" == "aarch64" ] ; then
    ARCH=arm64
else
    ARCH=amd64
fi

ls src -1 \
    | xargs -I {} basename {} .f \
    | xargs -I {} gfortran -fPIC -c -o /tmp/{}.o src/{}.f

ls src -1 \
    | xargs -I {} basename {} .f \
    | xargs -I {} echo /tmp/{}.o \
    | xargs ld -r -o /tmp/quadprog_${OS}_${ARCH}.syso

ls -alht /tmp

if [ "$OS" == "darwin" ]; then
    docker build . -f build-linux.dockerfile -t badgerodon-quadprog-build
    docker create --name extract badgerodon-quadprog-build:latest
    docker cp extract:/tmp/quadprog_linux_amd64.syso quadprog_linux_amd64.syso
    docker cp extract:/tmp/quadprog_linux_arm64.syso quadprog_linux_arm64.syso
    docker rm -f extract
fi
