FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

# add source code
ADD . /go/src/github.com/taka-wang/psmb

# add deps
WORKDIR /go/src/github.com/taka-wang/psmb
RUN glide install

# run test
RUN go test -v
