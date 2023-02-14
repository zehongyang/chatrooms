package models

import (
	"chat_room/db"
	"chat_room/proto/pb"
	"sync"
)

type DBRoom struct {
	*db.DbEngine
}

var GetDBRoom = func() func() *DBRoom {
	var (
		once sync.Once
		s    *DBRoom
	)
	return func() *DBRoom {
		once.Do(func() {
			s = &DBRoom{DbEngine: db.GetDBEngine()}
		})
		return s
	}
}()

func (s *DBRoom) ListRoom(q *pb.RoomListQuery) ([]*Room, error) {
	if q == nil {
		return nil, nil
	}
	var rooms []*Room
	err := s.Where("id > ?", q.Id).Limit(int(q.Limit), int(q.Skip)).Asc("id").Find(&rooms)
	return rooms, err
}
