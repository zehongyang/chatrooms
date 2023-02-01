package handlers

import (
	"chat_room/nets"
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
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		var res nets.ResponseInfo
		defer func() {
			if conn != nil {
				conn.WriteJSON(res)
			}
		}()
		if err != nil || conn == nil {
			logger.Error("WsConnectQuery", zap.Error(err))
			res.Code = http.StatusInternalServerError
			c.JSON(http.StatusBadRequest, res)
			return
		}
		var h nets.HeaderUser
		err = c.ShouldBindHeader(&h)
		if err != nil {
			logger.Error("WsConnectQuery", zap.Error(err))
			res.Code = http.StatusInternalServerError
			return
		}
		h.Uid = 1
		ss := nets.NewSession(conn, true, h.Uid)
		err = sm.Insert(ss)
		if err != nil {
			logger.Error("WsConnectQuery", zap.Error(err))
			res.Code = http.StatusInternalServerError
			return
		}
		res.Code = http.StatusOK
		go ss.Reader()
		go ss.Writer()
	}
}
