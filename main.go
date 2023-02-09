package main

import (
	"chat_room/handlers"
	"chat_room/nets"
	"github.com/gin-gonic/gin"
)

func main() {
	//注册处理websocket的方法
	nets.RegisterWebsocketHandler()
	engine := gin.Default()
	engine.GET("/ws", handlers.WsConnectQuery())
	engine.POST("/user/captcha", nets.WrapHttpFunc(handlers.CaptchaCreateQuery()))
	engine.POST("/user/register", nets.WrapHttpFunc(handlers.UserRegisterQuery()))
	engine.POST("/user/login", nets.WrapHttpFunc(handlers.UserLoginQuery()))
	engine.Run()
}
