FROM takawang/gozmq:x86
MAINTAINER Taka Wang <taka@cmwang.net>

# add source code
ADD . /go/src/github.com/taka-wang/psmb

# add deps
WORKDIR /go/src/github.com/taka-wang/psmb
RUN glide up

# run test
CMD ./test.sh