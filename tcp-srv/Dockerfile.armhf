FROM takawang/psmb:x86
MAINTAINER Taka Wang <taka@cmwang.net>

ADD . /go/
RUN cd /go/ \
    && go build -o psmbtcp \
    && cp psmbtcp /usr/bin/ 

CMD /usr/bin/psmbtcp

