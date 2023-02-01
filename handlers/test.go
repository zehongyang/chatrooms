package handlers

import (
	"chat_room/nets"
	"chat_room/proto/pb"
	"chat_room/utils/logger"
	"go.uber.org/zap"
	"net/http"
)

func TestHandler() nets.HandlerFunc {
	return func(ctx nets.ISessionContext) {
		var q pb.Hello
		var res pb.HelloResponse
		err := ctx.Request(&q)
		if err != nil {
			logger.Error("TestHandler", zap.Error(err))
			ctx.Response(http.StatusBadRequest, nil)
			return
		}
		res.Name = q.Name + "response"
		ctx.Response(http.StatusOK, &res)
	}
}
