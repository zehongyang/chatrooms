package nets

import (
	"chat_room/proto/pb"
	"chat_room/utils/errs"
	"chat_room/utils/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type ResponseInfo struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data proto.Message `json:"data"`
	Id   pb.HandlerId  `json:"id"`
}

type RequestInfo struct {
	Id   pb.HandlerId `json:"id"`
	Data string       `json:"data"`
}

type ISessionContext interface {
	Request(v proto.Message) error
	Response(code int, res proto.Message, msg ...string) error
	GetUid() int64
	GetIp() string
}

type HeaderUser struct {
	Uid   int64  `header:"Uid"`
	Token string `header:"Token"`
}

type HttpSessionContext struct {
	ctx    *gin.Context
	header HeaderUser
}

func (s HttpSessionContext) Request(v proto.Message) error {
	if s.ctx == nil {
		return errs.ErrorNilContext
	}
	err := s.ctx.ShouldBindHeader(&s.header)
	if err != nil {
		logger.Error("", zap.Any("err", err))
		return err
	}
	err = s.ctx.ShouldBind(v)
	if err != nil {
		logger.Error("", zap.Any("err", err))
		return err
	}
	return nil
}

func (s HttpSessionContext) Response(code int, res proto.Message, msgs ...string) error {
	var ri ResponseInfo
	ri.Code = code
	if len(msgs) > 0 {
		ri.Msg = msgs[0]
	}
	ri.Data = res
	s.ctx.JSON(http.StatusOK, ri)
	return nil
}

func (s HttpSessionContext) GetUid() int64 {
	return s.header.Uid
}

func (s HttpSessionContext) GetIp() string {
	return s.ctx.ClientIP()
}

type WebsocketSessionContext struct {
	session *Session
	data    []byte
	id      pb.HandlerId
}

func (s WebsocketSessionContext) Request(v proto.Message) error {
	return json.Unmarshal(s.data, v)
}

func (s WebsocketSessionContext) Response(code int, res proto.Message, msgs ...string) error {
	var ri ResponseInfo
	ri.Code = code
	if len(msgs) > 0 {
		ri.Msg = msgs[0]
	}
	ri.Data = res
	ri.Id = s.id
	data, err := json.Marshal(ri)
	if err != nil {
		return err
	}
	s.session.Write(data)
	return nil
}

func (s WebsocketSessionContext) GetUid() int64 {
	if s.session == nil {
		return 0
	}
	return s.session.Uid
}

func (s WebsocketSessionContext) GetIp() string {
	return s.session.LocalAddr().String()
}
