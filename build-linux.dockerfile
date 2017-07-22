FROM alpine:3.5 as amd64
RUN apk update
RUN apk add gfortran
WORKDIR /go/src/app
COPY ./src/ ./src/
COPY ./build.sh ./build.sh
RUN ./build.sh


FROM owlab/alpine-arm64:v3.5 as arm64
RUN apk update
RUN apk add gfortran
WORKDIR /go/src/app
COPY ./src/ ./src/
COPY ./build.sh ./build.sh
RUN ./build.sh


FROM alpine:3.5
COPY --from=amd64 /tmp/quadprog_linux_amd64.syso /tmp/quadprog_linux_amd64.syso
COPY --from=arm64 /tmp/quadprog_linux_arm64.syso /tmp/quadprog_linux_arm64.syso
