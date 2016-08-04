FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV PSMBTCP_CONFIG "/etc/psmbtcp"
#ENV SD_ENDPOINT "127.0.0.1:8500"
#ENV SD_BACKEND "consul"

WORKDIR /go
RUN go get -u github.com/taka-wang/psmb/tcp

ADD . /go/src/github.com/taka-wang/psmb/tcp-srv
# add config file
RUN mkdir -p ${PSMBTCP_CONFIG} && \ 
    cp /go/src/github.com/taka-wang/psmb/tcp/config.toml ${PSMBTCP_CONFIG}/

WORKDIR /go/src/github.com/taka-wang/psmb/tcp-srv
RUN go get -t ./... \
    && go build -o psmbtcp-srv \
    && cp psmbtcp-srv /usr/bin/ 

CMD /usr/bin/psmbtcp-srv

