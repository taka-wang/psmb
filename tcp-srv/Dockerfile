# psmbtcp-srv

FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV CONF_PSMBTCP "/etc/psmbtcp"
ENV EP_BACKEND "consul.cmwang.net:8500"

# add source code
ADD . /go/src/github.com/taka-wang/psmb
RUN ls

# install deps
RUN go get github.com/Masterminds/glide
RUN go build github.com/Masterminds/glide
WORKDIR /go/src/github.com/taka-wang/psmb/
RUN glide install

# add config file
RUN mkdir -p ${CONF_PSMBTCP} && \ 
    cp /go/src/github.com/taka-wang/psmb/tcp/config.toml ${CONF_PSMBTCP}/

# build
WORKDIR /go/src/github.com/taka-wang/psmb/tcp-srv
RUN go build -o psmbtcp-srv 
RUN cp psmbtcp-srv /usr/bin/ 

CMD /usr/bin/psmbtcp-srv

