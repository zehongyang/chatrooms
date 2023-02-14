package handlers

import (
	"chat_room/models"
	"chat_room/nets"
	"chat_room/proto/pb"
	"chat_room/utils/logger"
	"go.uber.org/zap"
	"net/http"
)

const (
	MaxSize = 100
)

func RoomListQuery() nets.HandlerFunc {
	dbRoom := models.GetDBRoom()
	return func(ctx nets.ISessionContext) {
		var (
			q   pb.RoomListQuery
			res pb.RoomListQueryResponse
		)
		err := ctx.Request(&q)
		if err != nil || q.Limit > MaxSize {
			logger.Error("RoomListQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusBadRequest, nil, "参数错误")
			return
		}
		rooms, err := dbRoom.ListRoom(&q)
		if err != nil {
			logger.Error("RoomListQuery", zap.Error(err), zap.Any("q", q))
			ctx.Response(http.StatusInternalServerError, nil)
			return
		}
		for _, room := range rooms {
			res.Rooms = append(res.Rooms, &pb.RoomInfo{
				Id:      room.Id,
				Name:    room.Name,
				Onlines: int32(room.Onlines),
				Img:     room.Img,
			})
		}
		ctx.Response(http.StatusOK, &res)
	}
}
