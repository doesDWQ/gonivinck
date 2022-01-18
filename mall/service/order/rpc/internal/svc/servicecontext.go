package svc

import (
	"mall/service/order/model"
	"mall/service/order/rpc/internal/config"
	"mall/service/product/rpc/productclient"
	"mall/service/user/rpc/userclient"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderModel model.OrderModel
	ProductRpc productclient.Product
	UserRpc    userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:     c,
		OrderModel: model.NewOrderModel(conn, c.CacheRedis),

		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpc)),

		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
