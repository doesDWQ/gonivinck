FROM golang:1.17

LABEL maintainer="Ving <ving@nivin.cn>"

ENV GOPROXY https://goproxy.cn,direct

# 安装必要的软件包和依赖包
USER root

# 安装依赖工具
RUN sed -i 's/deb.debian.org/mirrors.tuna.tsinghua.edu.cn/' /etc/apt/sources.list && \
    sed -i 's/security.debian.org/mirrors.tuna.tsinghua.edu.cn/' /etc/apt/sources.list && \
    sed -i 's/security-cdn.debian.org/mirrors.tuna.tsinghua.edu.cn/' /etc/apt/sources.list && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
    curl \
    zip \
    unzip \
    git \
    vim && \
    # 安装goctl工具
    GOPROXY=https://goproxy.cn/,direct go install github.com/tal-tech/go-zero/tools/goctl@cli && \
    # 安装protoc工具
    curl -L -o /tmp/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protoc-3.19.1-linux-x86_64.zip && \
    unzip -d /tmp/protoc /tmp/protoc.zip && \
    mv /tmp/protoc/bin/protoc $GOPATH/bin && \
    # 安装 protoc-gen-go
    go get -u github.com/golang/protobuf/protoc-gen-go@v1.4.0 && \
    # 清理垃圾
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    rm /var/log/lastlog /var/log/faillog

# $GOPATH/bin添加到环境变量中
ENV PATH $GOPATH/bin:$PATH

# 设置工作目录
WORKDIR /usr/src/code

CMD ["./main", "-f", "etc/config.yaml"]