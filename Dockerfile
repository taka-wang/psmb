#
# psmb
#
FROM takawang/ubuntu-gozmq
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
RUN go get github.com/marksalpeter/sugar && go get github.com/taka-wang/gocron


## Load app files
ADD . /go
RUN go build -o psmb

## Default command
#CMD ["go", "test", "-v"]
CMD ["/go/psmb"]
