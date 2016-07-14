#
# psmb
#
#FROM golang:1.6
FROM takawang/ubuntu-gozmq
MAINTAINER Taka Wang <taka@cmwang.net>

ADD .  /go/src/github.com/taka-wang/psmb

RUN echo "[url \"git@github.com:\"]\n\tinsteadOf = https://github.com/" >> /root/.gitconfig \
    && mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config \
    cd /go/src/github.com/taka-wang/psmb \
    && go get github.com/takawang/sugar \
    && go get github.com/taka-wang/gocron \
    && go get github.com/takawang/logrus 
    


#RUN go test -v
#RUN go build -o psmb

## Default command
CMD ["go", "test", "-v"]
#CMD ["/go/psmb"]

