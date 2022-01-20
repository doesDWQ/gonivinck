package logic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dtm-labs/dtmcli"
	"google.golang.org/grpc/codes"

	"mall/service/pay/model"
	"mall/service/pay/rpc/internal/svc"
	"mall/service/pay/rpc/pay"
	"mall/service/user/rpc/user"

	"github.com/dtm-labs/dtmgrpc"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
)

type CallbackRevertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallbackRevertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackRevertLogic {
	return &CallbackRevertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallbackRevertLogic) CallbackRevert(in *pay.CallbackRequest) (*pay.CallbackResponse, error) {
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

		res.Status = 0
		// 回滚更新支付状态
		err = l.svcCtx.PayModel.Update(res)
		if err != nil {
			return fmt.Errorf("回滚更新支付状态失败")
		}

		return nil
	}); err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &pay.CallbackResponse{}, nil
}
