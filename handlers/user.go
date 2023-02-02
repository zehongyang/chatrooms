package handlers

import (
	"chat_room/models"
	"chat_room/nets"
	"chat_room/proto/pb"
	"chat_room/utils/logger"
	"go.uber.org/zap"
	"net/http"
)

func UserRegisterQuery() nets.HandlerFunc {
	dbUser := models.GetDBUser()
	return func(ctx nets.ISessionContext) {
		var (
			q pb.UserRegisterQuery
			//res pb.UserRegisterQueryResponse
		)
		err := ctx.Request(&q)
		if err != nil || len(q.UserName) < 1 || len(q.Password) < 1 || len(q.IdKey) < 1 || len(q.Code) < 1 {
			logger.Error("UserRegisterQuery", zap.Error(err))
			ctx.Response(http.StatusBadRequest, nil)
			return
		}
		exists, err := dbUser.Exists(q.UserName)
		if err != nil {
			logger.Error("UserRegisterQuery", zap.Error(err))
			ctx.Response(http.StatusInternalServerError, nil)
			return
		}
		if exists {
			logger.Error("UserRegisterQuery", zap.Error(err))
			ctx.Response(http.StatusConflict, nil, "账号已被使用")
			return
		}
	}
}
