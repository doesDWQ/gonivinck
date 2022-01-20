package logic

import (
	"context"
	"database/sql"
	"fmt"
	"mall/service/pay/model"
	"mall/service/pay/rpc/internal/svc"
	"mall/service/pay/rpc/pay"
	"mall/service/user/rpc/user"

	"github.com/dtm-labs/dtmcli"
	"google.golang.org/grpc/codes"

	"github.com/dtm-labs/dtmgrpc"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
)

type CallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackLogic {
	return &CallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackLogic) Callback(in *pay.CallbackRequest) (*pay.CallbackResponse, error) {
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
		// 查询用户是否存在
		_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{
			Id: in.Uid,
		})
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		// 查询支付是否存在
		res, err := l.svcCtx.PayModel.FindOne(in.Id)
		if err != nil {
			if err == model.ErrNotFound {
				return fmt.Errorf("支付流水不存在")
			}
			return fmt.Errorf(err.Error())
		}
		// 支付金额与订单金额不符
		if in.Amount != res.Amount {
			return fmt.Errorf("支付金额与订单金额不符")
		}

		res.Source = in.Source
		res.Status = in.Status
		// 更新支付状态
		err = l.svcCtx.PayModel.Update(res)
		if err != nil {
			return fmt.Errorf("更新支付状态失败")
		}

		return nil
	}); err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &pay.CallbackResponse{}, nil
}
