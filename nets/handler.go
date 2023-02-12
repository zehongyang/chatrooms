package nets

import (
	"chat_room/proto/pb"
	"chat_room/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerFunc func(ctx ISessionContext)

func WrapHttpFunc(hf HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var h HeaderUser
		err := c.ShouldBindHeader(&h)
		if err != nil {
			logger.Error("WrapHttpFunc", zap.Error(err))
		}
		hf(&HttpSessionContext{
			ctx:    c,
			header: h,
		})
	}
}

var gWebsocketFunc = make(map[pb.HandlerId]HandlerFunc)

type WebsocketHandlerInfo struct {
	Id   pb.HandlerId
	Func HandlerFunc
}

func RegisterWebsocketHandler(his ...WebsocketHandlerInfo) {
	for _, hi := range his {
		if _, ok := gWebsocketFunc[hi.Id]; !ok {
			gWebsocketFunc[hi.Id] = hi.Func
		}
	}
}
