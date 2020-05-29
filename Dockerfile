FROM golang:latest

WORKDIR /root/workspace

RUN set -ex \
    && mkdir source && cd source \
    && git clone %#% \
    && cd django-go && go build -o django-go cmd/main.go

RUN set -ex \
    export GOPATH=~/git/middleware-contest-2020/gitlab

#go构建可执行文件
RUN set -ex \
    go build -o django-go cmd/main.go

#最终运行docker的命令
ENTRYPOINT  ["./django-go"]