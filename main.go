package main

import (
	"chat_room/handlers"
	"chat_room/nets"
	"chat_room/proto/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//注册处理websocket的方法
	initHandler()
	engine := gin.Default()
	engine.POST("/guest/captcha", nets.WrapHttpFunc(handlers.CaptchaCreateQuery()))
	engine.POST("/guest/register", nets.WrapHttpFunc(handlers.UserRegisterQuery()))
	engine.POST("/guest/login", nets.WrapHttpFunc(handlers.UserLoginQuery()))
	engine.StaticFS("/asset", http.Dir("C:/Users/Administrator/Desktop/tmp"))
	userGroup := engine.Group("/user")
	userGroup.Use(func(c *gin.Context) {
		//todo 鉴权认证
	})
	userGroup.GET("/ws", handlers.WsConnectQuery())
	userGroup.POST("/rooms", nets.WrapHttpFunc(handlers.RoomListQuery()))
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
