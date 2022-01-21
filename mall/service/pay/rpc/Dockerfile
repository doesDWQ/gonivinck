FROM registry.cn-shenzhen.aliyuncs.com/dwq/go-zero:v1.0

WORKDIR /usr/src/code
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY service/pay/rpc/etc etc
RUN go build -ldflags="-s -w" -o pay service/pay/rpc/pay.go

CMD ["./pay", "-f", "etc/pay.yaml"]