package nets

import (
	"chat_room/proto/pb"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx ISessionContext)

func WrapHttpFunc(hf HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		hf(&HttpSessionContext{
			ctx:    c,
			header: HeaderUser{},
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
