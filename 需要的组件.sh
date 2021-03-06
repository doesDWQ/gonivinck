需要的组件
etcd                1
grafana             1
jaeger              1
mysql               1
prometheus          1
redis               1


docker run --name prometheus -d \
    -p 9090:9090 \
    -v ./prometheus:/etc/prometheus \
    prom/prometheus

引入bazel加快编译速度
更新依赖
bazel run //:gazelle -- update-repos -from_file=go.mod 
bazel run //:gazelle
编译所有：






docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.30

  docker run -d -p 3000:3000 --name grafana grafana/grafana-enterprise:8.2.0


docker run -d --name dtm \
-p 36789:36789 \
-p 36790:36790  \
-v ${PWD}/dtm/config:/app/dtm/configs \
registry.cn-shenzhen.aliyuncs.com/dwq/dtm:v1.0

集成的swagger文档
https://github.com/zeromicro/goctl-swagger