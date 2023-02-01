package handlers

import (
	"chat_room/nets"
	"chat_room/proto/pb"
	"encoding/base64"
	"fmt"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

func CaptchaCreateQuery() nets.HandlerFunc {
	var configD = base64Captcha.ConfigDigit{
		Height:     40,
		Width:      100,
		MaxSkew:    0.5,
		DotCount:   50,
		CaptchaLen: 4,
	}
	return func(ctx nets.ISessionContext) {
		var res pb.CaptchaCreateQueryResponse
		id, instance := base64Captcha.GenerateCaptcha("", configD)
		res.Id = id
		res.Bs4 = fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(instance.BinaryEncoding()))
		ctx.Response(http.StatusOK, &res)
	}
}
