FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV PSMBTCP_CONFIG "/etc/psmbtcp"
#ENV CONSUL_ENDPOINT "127.0.0.1:8500"

WORKDIR /go
RUN go get -u github.com/taka-wang/psmb
RUN go get -u github.com/taka-wang/psmb/tcp
RUN go get -u github.com/spf13/viper
#RUN go get github.com/spf13/viper/remote

## Load app files
ADD . /go/src/github.com/taka-wang/psmb/redis-writer
# add config file
ADD config.toml ${PSMBTCP_CONFIG}/
WORKDIR /go/src/github.com/taka-wang/psmb/redis-writer

RUN go get -t ./...

## Default command
CMD ["go", "test", "-v"]