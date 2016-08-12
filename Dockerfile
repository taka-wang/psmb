FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

# add source code
ADD . /go/src/github.com/taka-wang/psmb

# add deps
WORKDIR /go/src/github.com/taka-wang/psmb
RUN glide install
RUN apt-get update && apt-get install curl -y

# run test
RUN go test -v
RUN bash <(curl -s https://codecov.io/bash) -t 558aa53d-c58d-4df4-a1c1-a22a6e6d8572
