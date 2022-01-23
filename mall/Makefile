productRpc:
	cd service/product/rpc && go run product.go

userRpc:
	cd service/user/rpc && go run user.go


orderRpc:
	cd service/order/rpc && go run order.go

payRpc:
	cd service/pay/rpc && go run pay.go

orderApi:
	cd service/order/api && go run order.go

productApi:
	cd service/product/api && go run product.go

userApi:
	cd service/user/api && go run user.go

payApi:
	cd service/pay/api && go run pay.go


CURRENT_DIR := ${PWD}
# 编译所有的文件
build:
	cd service/product/rpc/ && go build -o product.exe product.go
	cd service/product/api/ && go build -o product.exe product.go

	cd service/user/rpc/ && go build -o user.exe user.go
	cd service/user/api/ && go build -o user.exe user.go

	cd service/order/rpc/ && go build -o order.exe order.go
	cd service/order/api/ && go build -o order.exe order.go

	cd service/pay/rpc/ && go build -o pay.exe pay.go
	cd service/pay/api/ && go build -o pay.exe pay.go
