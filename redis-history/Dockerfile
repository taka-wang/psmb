# redis-history

FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV CONF_PSMBTCP "/etc/psmbtcp"
ENV EP_BACKEND "consul.cmwang.net:8500"

WORKDIR /go
RUN go get github.com/taka-wang/psmb/viper-conf
RUN go get -u github.com/taka-wang/psmb/tcp
RUN go get github.com/taka-wang/psmb


## Load app files
ADD . /go/src/github.com/taka-wang/psmb/redis-history
# add config file
RUN mkdir -p ${CONF_PSMBTCP} && \ 
    cp /go/src/github.com/taka-wang/psmb/tcp/config.toml ${CONF_PSMBTCP}/

WORKDIR /go/src/github.com/taka-wang/psmb/redis-history

RUN go get -t ./...

## Default command
CMD ./test.sh