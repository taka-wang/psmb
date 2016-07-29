FROM takawang/gozmq:x86
MAINTAINER Taka Wang <taka@cmwang.net>

ADD . /go/src/github.com/taka-wang/psmb
RUN go get github.com/takawang/sugar \
    && go get github.com/taka-wang/gocron \
    && go get github.com/takawang/logrus \ 
    && cd /go/src/github.com/taka-wang/psmb \
    && go test -v

#CMD go test -v