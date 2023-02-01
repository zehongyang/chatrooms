package main

import (
	"chat_room/proto/pb"
	"encoding/base64"
	"github.com/mojocn/base64Captcha"
	"google.golang.org/protobuf/proto"
	"os"
	"testing"
)

func TestS(t *testing.T) {
	t.Log(-3 % 2)
	var tss pb.Hello
	_, err := proto.Marshal(&tss)
	t.Log(err)
}

func TestCaptcha(t *testing.T) {
	var configD = base64Captcha.ConfigDigit{
		Height:     40,
		Width:      100,
		MaxSkew:    0.5,
		DotCount:   50,
		CaptchaLen: 4,
	}
	_, instance := base64Captcha.GenerateCaptcha("", configD)
	toString := base64.StdEncoding.EncodeToString(instance.BinaryEncoding())
	t.Log(toString)
	os.WriteFile("./ss.png", instance.BinaryEncoding(), 0755)
}
