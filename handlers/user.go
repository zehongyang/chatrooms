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
			logger.Error("UserRegisterQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusBadRequest, nil, "参数错误")
			return
		}
		if len(q.UserName) < 6 || len(q.UserName) > 16 {
			logger.Error("UserRegisterQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusBadRequest, nil, "账号长度为6-16位")
			return
		}
		if len(q.Password) < 6 {
			logger.Error("UserRegisterQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusBadRequest, nil, "密码长度为6位及以上")
			return
		}
		//验证码是否正确
		if !base64Captcha.VerifyCaptcha(q.IdKey, q.Code) {
			logger.Error("UserRegisterQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusBadRequest, nil, "验证码错误")
			return
		}
		exists, err := dbUser.Exists(q.UserName)
		if err != nil {
			logger.Error("UserRegisterQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusInternalServerError, nil)
			return
		}
		if exists {
			logger.Error("UserRegisterQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusConflict, nil, "账号已被使用")
			return
		}
		u, err := dbUser.Register(&q, ctx.GetIp())
		if err != nil || u == nil {
			logger.Error("UserRegisterQuery", zap.Error(err), zap.Any("q", q))
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

func UserLoginQuery() nets.HandlerFunc {
	dbUser := models.GetDBUser()
	return func(ctx nets.ISessionContext) {
		var (
			q   pb.UserLoginQuery
			res pb.UserLoginQueryResponse
		)
		err := ctx.Request(&q)
		if err != nil || len(q.UserName) < 1 || len(q.Password) < 1 || len(q.IdKey) < 1 || len(q.Code) < 1 {
			logger.Error("UserLoginQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusBadRequest, nil, "参数错误")
			return
		}
		//验证码是否正确
		if !base64Captcha.VerifyCaptcha(q.IdKey, q.Code) {
			logger.Error("UserLoginQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusBadRequest, nil, "验证码错误")
			return
		}
		u, err := dbUser.Login(q.UserName, q.Password, ctx.GetIp())
		if err != nil {
			logger.Error("UserRegisterQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusInternalServerError, nil)
			return
		}
		if u == nil {
			logger.Error("UserLoginQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusBadRequest, nil, "用户名或密码错误")
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
