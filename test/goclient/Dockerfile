#
# golang client to modbusd
#
FROM takawang/dummy-psmb
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
#RUN go get github.com/takawang/sugar

## Load app files
ADD . /go

## Default command
CMD ["go", "test", "-v"]

