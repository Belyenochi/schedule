FROM registry.cn-hangzhou.aliyuncs.com/middlewarerace2020/go

WORKDIR /root/workspace

ENV GOPATH /root/workspace

RUN set -ex \
    && mkdir src && cd src && git clone https://code.aliyun.com/middleware-contest-2020/django-go.git \
    && cd django-go && ls -l  && /usr/local/go/bin/go build -o django.go cmd/main.go

#最终运行docker的命令
ENTRYPOINT  ["/root/workspace/src/django-go/django.go"]