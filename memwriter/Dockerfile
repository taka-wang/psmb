FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
#RUN go get -u github.com/taka-wang/psmb

## Load app files
ADD . /go/src/github.com/taka-wang/psmb/memwriter
WORKDIR /go/src/github.com/taka-wang/psmb/memwriter
RUN go get -t ./... && go test -v