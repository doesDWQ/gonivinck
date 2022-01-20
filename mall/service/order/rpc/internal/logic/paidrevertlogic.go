package logic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dtm-labs/dtmcli"
	"google.golang.org/grpc/codes"

	"mall/service/order/model"
	"mall/service/order/rpc/internal/svc"
	"mall/service/order/rpc/order"

	"github.com/dtm-labs/dtmgrpc"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
)

type PaidRevertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPaidRevertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaidRevertLogic {
	return &PaidRevertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PaidRevertLogic) PaidRevert(in *order.PaidRequest) (*order.PaidResponse, error) {
	// 获取 RawDB
	db, err := sqlx.NewMysql(l.svcCtx.Config.Mysql.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	// 获取 barrier，用于防止空补偿、空悬挂
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 查询订单是否存在
		res, err := l.svcCtx.OrderModel.FindOne(in.Id)
		if err != nil {
			if err == model.ErrNotFound {
				return fmt.Errorf("订单不存在")
			}
			return fmt.Errorf(err.Error())
		}

		res.Status = 0
		// 回滚更新订单状态
		err = l.svcCtx.OrderModel.Update(res)
		if err != nil {
			return fmt.Errorf("回滚更新订单状态失败")
		}

		return nil
	}); err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &order.PaidResponse{}, nil
}
