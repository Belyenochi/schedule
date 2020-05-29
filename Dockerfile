FROM centos:centos7

WORKDIR /root/workspace

RUN set -ex  && mkdir src && cd src

RUN set -ex && yum install -y wget && yum install -y git

## tar zxvf
RUN set -ex  && wget -q  https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz && tar zxf *.tar.gz -C /usr/local/

ENV GOROOT /usr/local/go

ENV PATH $PATH:$GOROOT/bin

ENV GOPATH /root/workspace

RUN set -ex \
    && cd src && git clone https://code.aliyun.com/middleware-contest-2020/django-go.git \
    && cd django-go && ls -l  && /usr/local/go/bin/go build -o django.go cmd/main.go

#最终运行docker的命令
ENTRYPOINT  ["/root/workspace/src/django-go/django.go"]