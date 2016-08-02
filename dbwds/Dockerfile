FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
RUN go get -u github.com/taka-wang/psmb && go get -u github.com/taka-wang/psmb/tcp

## Load app files
ADD . /go/src/github.com/taka-wang/psmb/dbwds
WORKDIR /go/src/github.com/taka-wang/psmb/dbwds

RUN go get -t ./...

## Default command
CMD ["go", "test", "-v"]