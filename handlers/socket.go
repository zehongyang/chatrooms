package handlers

import (
	"chat_room/nets"
	"chat_room/proto/pb"
	"chat_room/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

func WsConnectQuery() gin.HandlerFunc {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	sm := nets.GetSessionManager()
	return func(c *gin.Context) {
		var res nets.ResponseInfo
		var h nets.HeaderUser
		err := c.ShouldBindHeader(&h)
		if err != nil || h.Uid < 1 {
			logger.Error("WsConnectQuery", zap.Error(err))
			res.Code = http.StatusInternalServerError
			c.JSON(http.StatusBadRequest, res)
			return
		}
		logger.Info("", zap.Any("h", h))
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil || conn == nil {
			logger.Error("WsConnectQuery", zap.Error(err))
			res.Code = http.StatusInternalServerError
			c.JSON(http.StatusBadRequest, res)
			return
		}
		ss := nets.NewSession(conn, true, h.Uid)
		err = sm.Insert(ss)
		if err != nil {
			logger.Error("WsConnectQuery", zap.Error(err))
			res.Code = http.StatusInternalServerError
			c.JSON(http.StatusBadRequest, res)
			return
		}
		res.Code = http.StatusOK
		go ss.Reader()
		go ss.Writer()
	}
}

const (
	Ping = "ping"
	Pong = "pong"
)

func HeartBeatQuery() nets.HandlerFunc {
	return func(ctx nets.ISessionContext) {
		var (
			q   pb.HeartBeatQuery
			res pb.HeartBeatQueryResponse
		)
		uid := ctx.GetUid()
		if uid < 1 {
			ctx.Response(http.StatusBadRequest, nil)
			return
		}
		err := ctx.Request(&q)
		if err != nil {
			logger.Error("HeartBeatQuery", zap.Error(err), zap.Any("q", &q))
			ctx.Response(http.StatusBadRequest, nil)
			return
		}
		if q.Msg != Ping {
			logger.Error("HeartBeatQuery", zap.Any("q", &q))
			ctx.Response(http.StatusBadRequest, nil)
			return
		}
		res.Msg = Pong
		ctx.Response(http.StatusOK, &res)
	}
}
