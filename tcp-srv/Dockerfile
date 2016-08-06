# psmbtcp-srv

FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

ENV CONF_PSMBTCP "/etc/psmbtcp"
ENV EP_BACKEND "consul.cmwang.net:8500"

WORKDIR /go

RUN go get github.com/spf13/viper
RUN go get github.com/takawang/logrus
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/garyburd/redigo/redis

RUN go get -u github.com/taka-wang/psmb/tcp

ADD . /go/src/github.com/taka-wang/psmb/tcp-srv
# add config file
RUN mkdir -p ${CONF_PSMBTCP} && \ 
    cp /go/src/github.com/taka-wang/psmb/tcp/config.toml ${CONF_PSMBTCP}/

WORKDIR /go/src/github.com/taka-wang/psmb/tcp-srv
#RUN go get -t ./... 
RUN go build -o psmbtcp-srv 
RUN cp psmbtcp-srv /usr/bin/ 

CMD /usr/bin/psmbtcp-srv

