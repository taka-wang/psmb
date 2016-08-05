# mbtcp

FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV CONF_PSMBTCP "/etc/psmbtcp"
ENV EP_BACKEND "dev.cmwang.net:8500"

WORKDIR /go
RUN go get -u github.com/spf13/viper
#RUN go get github.com/spf13/viper/remote
RUN go get -u github.com/taka-wang/psmb

## Load app files
ADD . /go/src/github.com/taka-wang/psmb/tcp
ADD config.toml ${CONF_PSMBTCP}

WORKDIR /go/src/github.com/taka-wang/psmb/tcp

RUN go get -t ./... && go test -v
