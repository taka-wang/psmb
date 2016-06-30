#
# psmb
#
FROM takawang/ubuntu-gozmq
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
RUN go get github.com/marksalpeter/sugar


## Load app files
ADD . /go
RUN go build

## Default command
#CMD ["go", "test", "-v"]
CMD ["/go/psmb"]
