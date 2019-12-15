FROM golang:buster AS build-env

LABEL maintainer=alandiegosantos@gmail.com

RUN apt-get update && apt-get install -y protobuf-compiler

COPY . $GOPATH/src/http-random-stress

WORKDIR $GOPATH/src/http-random-stress

RUN make -f build/Makefile

FROM gcr.io/distroless/base

COPY --from=build-env /go/src/http-random-stress/http-random-stress  /usr/bin/http-random-stress

ENTRYPOINT ["http-random-stress server --crt /app/assets/server.crt --key /app/assets/server.key"]