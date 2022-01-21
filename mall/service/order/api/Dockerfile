FROM registry.cn-shenzhen.aliyuncs.com/dwq/go-zero:v1.0

WORKDIR /usr/src/code
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY service/order/api/etc etc
RUN go build -ldflags="-s -w" -o order service/order/api/order.go

CMD ["./order", "-f", "etc/order.yaml"]