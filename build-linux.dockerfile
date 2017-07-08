FROM alpine:3.6

RUN apk update
RUN apk add gfortran

WORKDIR /go/src/app
COPY . .

RUN ./build.sh
