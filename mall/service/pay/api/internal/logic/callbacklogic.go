package logic

import (
	"context"

	"mall/service/order/rpc/order"
	"mall/service/pay/api/internal/svc"
	"mall/service/pay/api/internal/types"
	"mall/service/pay/rpc/pay"

	"github.com/dtm-labs/dtmgrpc"
	"github.com/tal-tech/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type CallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) CallbackLogic {
	return CallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CallbackLogic) Callback(req types.CallbackRequest) (resp *types.CallbackResponse, err error) {
	// 获取 PayRpc BuildTarget
	payRpcBusiServer, err := l.svcCtx.Config.PayRpc.BuildTarget()
	if err != nil {
		return nil, status.Error(100, "支付回调异常")
	}
	// 获取 OrderRpc BuildTarget
	orderRpcBusiServer, err := l.svcCtx.Config.OrderRpc.BuildTarget()
	if err != nil {
		return nil, status.Error(100, "支付回调异常")
	}
	// dtm 服务的 etcd 注册地址
	var dtmServer = "etcd://9.135.226.207:2379/dtmservice"
	// 创建一个gid
	gid := dtmgrpc.MustGenGid(dtmServer)
	// 创建一个saga协议的事务
	saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
		Add(payRpcBusiServer+"/payclient.Pay/Callback", payRpcBusiServer+"/payclient.Pay/CallbackRevert",
			&pay.CallbackRequest{
				Id:     req.Id,
				Uid:    req.Uid,
				Oid:    req.Oid,
				Amount: req.Amount,
				Source: req.Source,
				Status: req.Status,
			}).
		Add(orderRpcBusiServer+"/orderclient.Order/Paid", orderRpcBusiServer+"/orderclient.Order/PaidRevert", &order.PaidRequest{
			Id: req.Oid,
		})

	// 事务提交
	err = saga.Submit()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	return &types.CallbackResponse{}, nil
}
