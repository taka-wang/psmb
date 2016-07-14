#
# psmb
#
FROM takawang/ubuntu-gozmq
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
RUN go get github.com/takawang/sugar \
    && go get github.com/taka-wang/gocron \
    && go get github.com/takawang/logrus \
    && git config --global url."git@github.com:".insteadOf "https://github.com/" \
    && go get github.com/taka-wang/psmb 


## Load app files
ADD . /go
#RUN go test -v
#RUN go build -o psmb

## Default command
CMD ["go", "test", "-v"]
#CMD ["/go/psmb"]
