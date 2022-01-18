package handler

import (
	"net/http"

	"mall/service/product/api/internal/logic"
	"mall/service/product/api/internal/svc"
	"mall/service/product/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func CreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logx.Info(r.Form)
		var req types.CreateRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Info("解析出错了")
			httpx.Error(w, err)
			return
		}

		l := logic.NewCreateLogic(r.Context(), svcCtx)
		resp, err := l.Create(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
