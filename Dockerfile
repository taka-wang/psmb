#
# psmb
#
#FROM golang:1.6
FROM takawang/ubuntu-gozmq
MAINTAINER Taka Wang <taka@cmwang.net>

ADD . /go/src/github.com/taka-wang/psmb
RUN go get github.com/takawang/sugar \
    && go get github.com/taka-wang/gocron \
    && go get github.com/takawang/logrus \ 
    && cd /go/src/github.com/taka-wang/psmb \
    && go test -v \
    && go build \
	&& cd / \
    && git clone https://github.com/taka-wang/psmb-srv.git \
    && cd psmb-srv \
    && go build -o psmb \
    && cp psmb /usr/bin/ 

CMD /usr/bin/psmb

