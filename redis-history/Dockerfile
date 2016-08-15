# redis-history

FROM takawang/gozmq:x86
MAINTAINER Taka Wang <taka@cmwang.net>

ENV CONF_PSMBTCP "/etc/psmbtcp"
ENV EP_BACKEND "consul.cmwang.net:8500"

# add source code from root
ADD . /go/src/github.com/taka-wang/psmb

# install deps
WORKDIR /go/src/github.com/taka-wang/psmb/
RUN glide up

# add config file
RUN mkdir -p ${CONF_PSMBTCP} && \ 
    cp /go/src/github.com/taka-wang/psmb/tcp/config.toml ${CONF_PSMBTCP}/

WORKDIR /go/src/github.com/taka-wang/psmb/redis-history

## Default command
CMD ./test.sh