package main

import (
	"chat_room/handlers"
	"chat_room/nets"
	"chat_room/proto/pb"
	"chat_room/utils/logger"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	//注册处理websocket的方法
	initHandler()
	engine := gin.Default()
	engine.POST("/guest/captcha", nets.WrapHttpFunc(handlers.CaptchaCreateQuery()))
	engine.POST("/guest/register", nets.WrapHttpFunc(handlers.UserRegisterQuery()))
	engine.POST("/guest/login", nets.WrapHttpFunc(handlers.UserLoginQuery()))
	userGroup := engine.Group("/user")
	userGroup.Use(func(c *gin.Context) {
		//todo 鉴权认证
	})
	userGroup.GET("/ws", handlers.WsConnectQuery())
	go func() {
		sm := nets.GetSessionManager()
		for {
			ss, err := sm.GetSession(21)
			if err != nil {
				logger.Error("", zap.Error(err))
			}
			if ss != nil {
				var res nets.ResponseInfo
				res.Code = http.StatusOK
				res.Data = &pb.ChatMessage{
					Content: fmt.Sprintf("消息%s", time.Now().Format("2006-01-02 15:04:05")),
					MsgType: pb.MsgType_MT_Text,
					RoomId:  1,
				}
				res.Id = pb.HandlerId_HI_Message
				marshal, err := json.Marshal(res)
				if err != nil {
					logger.Error("", zap.Error(err))
					continue
				}
				ss.Write(marshal)
			}
			time.Sleep(time.Second * 5)
		}
	}()
	engine.Run()
}

func initHandler() {
	nets.RegisterWebsocketHandler([]nets.WebsocketHandlerInfo{
		{
			Id:   pb.HandlerId_HI_HeartBeat,
			Func: handlers.HeartBeatQuery(),
		},
	}...)
}
