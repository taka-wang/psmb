FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
RUN go get -u github.com/taka-wang/psmb
RUN go get -u github.com/taka-wang/psmb/tcp
RUN go get -u github.com/spf13/viper
RUN go get github.com/spf13/viper/remote

## Load app files
ADD . /go/src/github.com/taka-wang/psmb/redis-history
WORKDIR /go/src/github.com/taka-wang/psmb/redis-history

RUN go get -t ./...

## Default command
CMD ["go", "test", "-v"]