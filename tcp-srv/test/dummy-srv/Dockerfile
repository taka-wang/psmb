FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV CONF_PSMBTCP "/etc/psmbtcp"
ENV EP_BACKEND "consul.cmwang.net:8500"

# add source code from root
ADD . /go/src/github.com/taka-wang/psmb

# install deps
WORKDIR /go/src/github.com/taka-wang/psmb/
RUN glide install
RUN apt-get update && apt-get install curl -y

# set work dir
WORKDIR /go/src/github.com/taka-wang/psmb/tcp-srv/test/dummy-srv

CMD ./test.sh
