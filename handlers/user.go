package handlers

import (
	"chat_room/models"
	"chat_room/nets"
	"chat_room/proto/pb"
	"chat_room/utils/logger"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"time"
)

func UserRegisterQuery() nets.HandlerFunc {
	dbUser := models.GetDBUser()
	rand.Seed(time.Now().UnixNano())
	return func(ctx nets.ISessionContext) {
		var (
			q   pb.UserRegisterQuery
			res pb.UserRegisterQueryResponse
		)
		err := ctx.Request(&q)
		if err != nil || len(q.UserName) < 1 || len(q.Password) < 1 || len(q.IdKey) < 1 || len(q.Code) < 1 {
			logger.Error("UserRegisterQuery", zap.Error(err))
			ctx.Response(http.StatusBadRequest, nil)
			return
		}
		//验证码是否正确
		if !base64Captcha.VerifyCaptcha(q.IdKey, q.Code) {
			logger.Error("UserRegisterQuery", zap.Error(err))
			ctx.Response(http.StatusBadRequest, nil, "验证码错误")
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
		u, err := dbUser.Register(&q, ctx.GetIp())
		if err != nil || u == nil {
			logger.Error("UserRegisterQuery", zap.Error(err))
			ctx.Response(http.StatusInternalServerError, nil)
			return
		}
		res.UserInfo = &pb.UserInfo{
			Uid:      u.Id,
			NickName: u.NickName,
			Avatar:   u.Avatar,
			Token:    u.Token,
		}
		ctx.Response(http.StatusOK, &res)
	}
}
