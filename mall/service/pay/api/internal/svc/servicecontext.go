package svc

import (
	"mall/service/order/rpc/orderclient"
	"mall/service/pay/api/internal/config"
	"mall/service/pay/rpc/payclient"

	"github.com/tal-tech/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	PayRpc   payclient.Pay
	OrderRpc orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		PayRpc:   payclient.NewPay(zrpc.MustNewClient(c.PayRpc)),
		OrderRpc: orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
	}
}
