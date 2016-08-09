# redis-filter

FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV CONF_PSMBTCP "/etc/psmbtcp"
ENV EP_BACKEND "consul.cmwang.net:8500"

WORKDIR /go
RUN go get -u github.com/taka-wang/psmb

## Load app files
ADD . /go/src/github.com/taka-wang/psmb/redis-filter
# add config file
RUN mkdir -p ${CONF_PSMBTCP} && \ 
    cp /go/src/github.com/taka-wang/psmb/tcp/config.toml ${CONF_PSMBTCP}/

WORKDIR /go/src/github.com/taka-wang/psmb/redis-filter
RUN go get -t ./... 

## Default command
CMD ./test.sh