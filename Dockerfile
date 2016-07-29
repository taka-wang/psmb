FROM takawang/gozmq:x86
MAINTAINER Taka Wang <taka@cmwang.net>

ADD . /go/src/github.com/taka-wang/psmb
WORKDIR /go/src/github.com/taka-wang/psmb
RUN go go get -t ./... && go test -v
