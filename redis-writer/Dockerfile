FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

#ENV CONF_PSMBTCP "/etc/psmbtcp"
ENV CONF_PSMBTCP "/psmbtcp/config.toml"
ENV EP_BACKEND "dev.cmwang.net:8500"

WORKDIR /go
RUN go get -u github.com/spf13/viper
#RUN go get github.com/spf13/viper/remote
RUN go get -u github.com/taka-wang/psmb
RUN go get -u github.com/taka-wang/psmb/tcp

## Load app files
ADD . /go/src/github.com/taka-wang/psmb/redis-writer
# add config file
RUN mkdir -p ${PSMBTCP_CONFIG} && \ 
    cp /go/src/github.com/taka-wang/psmb/tcp/config.toml ${PSMBTCP_CONFIG}/

WORKDIR /go/src/github.com/taka-wang/psmb/redis-writer

RUN go get -t ./...

## Default command
CMD ["go", "test", "-v"]