FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV PSMBTCP_CONFIG "/etc/psmbtcp"
#ENV SD_ENDPOINT "127.0.0.1:8500"
#ENV SD_BACKEND "consul"

WORKDIR /go
RUN go get -u github.com/spf13/viper
#RUN go get github.com/spf13/viper/remote
RUN go get -u github.com/taka-wang/psmb

## Load app files
ADD . /go/src/github.com/taka-wang/psmb/tcp
ADD config.toml ${PSMBTCP_CONFIG}

WORKDIR /go/src/github.com/taka-wang/psmb/tcp

RUN go get -t ./... && go test -v
