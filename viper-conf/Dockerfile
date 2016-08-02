# mbtcp

FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV CONF_PSMBTCP "/etc/psmbtcp"
ENV EP_BACKEND "consul.cmwang.net:8500"

WORKDIR /go

RUN go get -u github.com/taka-wang/psmb/tcp
RUN go get -u github.com/spf13/viper
RUN go get -u github.com/apex/log

## Load app files
ADD . /go/src/github.com/taka-wang/psmb/viper-conf
# add config file
RUN mkdir -p ${CONF_PSMBTCP} && \ 
    cp /go/src/github.com/taka-wang/psmb/tcp/config.toml ${CONF_PSMBTCP}/

WORKDIR /go/src/github.com/taka-wang/psmb/viper-conf

RUN go get -t ./... 
RUN go test -v
