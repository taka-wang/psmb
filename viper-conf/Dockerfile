# mbtcp

FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV CONF_PSMBTCP "/etc/psmbtcp"
ENV EP_BACKEND "consul.cmwang.net:8500"

WORKDIR /go


## Load app files
ADD . /go/src/github.com/taka-wang/psmb/viper-conf
ADD config.toml ${CONF_PSMBTCP}

WORKDIR /go/src/github.com/taka-wang/psmb/viper-conf

RUN go get -t ./... 
RUN go test -v
