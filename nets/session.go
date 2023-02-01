package nets

import (
	"chat_room/proto/pb"
	"chat_room/utils/errs"
	"chat_room/utils/logger"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

const (
	ChanNums     = 1024
	ReadTimeOut  = time.Second * 50
	WriteTimeOut = time.Second * 10
)

type Session struct {
	*websocket.Conn
	IsLogin   bool
	Uid       int64
	sm        *SessionManager
	Closed    bool
	writeChan chan []byte
	ctx       context.Context
}

func NewSession(conn *websocket.Conn, isLogin bool, uid int64) *Session {
	return &Session{
		Conn:      conn,
		IsLogin:   isLogin,
		Uid:       uid,
		sm:        GetSessionManager(),
		Closed:    false,
		writeChan: make(chan []byte, 1024),
		ctx:       context.Background(),
	}
}

func (s *Session) CloseSession() error {
	if s.sm == nil {
		return errs.ErrorNilSessionManager
	}
	if s.Closed {
		return nil
	}
	err := s.sm.Close(s.Uid)
	if err != nil {
		logger.Error("", zap.Error(err), zap.Any("s", s))
		return err
	}
	return nil
}

func (s *Session) Reader() {
	defer func() {
		s.CloseSession()
	}()
	for {
		err := s.SetReadDeadline(time.Now().Add(ReadTimeOut))
		if err != nil {
			logger.Error("", zap.Error(err), zap.Any("s", s))
			return
		}
		_, data, err := s.ReadMessage()
		if err != nil {
			logger.Error("", zap.Error(err), zap.Any("s", s))
			return
		}
		s.dealRead(data)
	}
}

func (s *Session) dealRead(data []byte) {
	var ri RequestInfo
	var res ResponseInfo
	err := json.Unmarshal(data, &ri)
	if err != nil {
		logger.Error("", zap.Error(err))
		res.Code = http.StatusBadRequest
		resData, err := json.Marshal(res)
		if err != nil {
			logger.Error("", zap.Error(err))
		} else {
			s.Write(resData)
		}
		return
	}
	if !s.IsLogin {
		s.CloseSession()
		return
	}
	sctx := WebsocketSessionContext{
		session: s,
		data:    []byte(ri.Data),
	}
	h, ok := s.sm.funcMap[ri.Id]
	if !ok {
		res.Code = http.StatusNotFound
		resData, err := json.Marshal(res)
		if err != nil {
			logger.Error("", zap.Error(err))
		} else {
			s.Write(resData)
		}
		return
	}
	h(sctx)
}

func (s *Session) Write(data []byte) {
	s.writeChan <- data
}

func (s *Session) Writer() {
	defer func() {
		s.CloseSession()
	}()
	for {
		select {
		case data, ok := <-s.writeChan:
			if !ok {
				return
			}
			err := s.SetWriteDeadline(time.Now().Add(WriteTimeOut))
			if err != nil {
				logger.Error("", zap.Error(err), zap.Any("s", s))
			}
			err = s.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logger.Error("", zap.Error(err), zap.Any("s", s))
				return
			}
		}
	}
}

type SessionMap struct {
	lock sync.RWMutex
	sm   map[int64]*Session
}

type sessionSlice struct {
	buckets []SessionMap
	length  int64
}

func (s *sessionSlice) WithBuckets(n int64) *sessionSlice {
	s.buckets = make([]SessionMap, n)
	s.length = n
	return s
}

func NewSessionSlice() *sessionSlice {
	length := 4
	return &sessionSlice{buckets: make([]SessionMap, length), length: int64(length)}
}

type SessionManager struct {
	sessions *sessionSlice
	funcMap  map[pb.HandlerId]HandlerFunc
}

var GetSessionManager = func() func() *SessionManager {
	var (
		once sync.Once
		s    *SessionManager
	)
	return func() *SessionManager {
		once.Do(func() {
			s = &SessionManager{
				sessions: NewSessionSlice().WithBuckets(16),
				funcMap:  gWebsocketFunc,
			}
		})
		return s
	}
}()

func (s *SessionManager) Insert(ss *Session) error {
	if ss == nil || ss.Conn == nil {
		return errs.ErrorNilConn
	}
	if ss.Uid < 1 {
		return errs.ErrorInvalidUid
	}
	idx := ss.Uid % s.sessions.length
	if idx >= s.sessions.length {
		return errs.ErrorInvalidSessionSlice
	}
	if ss.sm == nil {
		ss.sm = s
	}
	s.sessions.buckets[idx].lock.Lock()
	defer s.sessions.buckets[idx].lock.Unlock()
	if s.sessions.buckets[idx].sm == nil {
		s.sessions.buckets[idx].sm = make(map[int64]*Session)
	}
	if oldSs, ok := s.sessions.buckets[idx].sm[ss.Uid]; ok && !oldSs.Closed {
		oldSs.Close()
		oldSs.Closed = true
		close(oldSs.writeChan)
	}
	s.sessions.buckets[idx].sm[ss.Uid] = ss
	return nil
}

func (s *SessionManager) Close(uid int64) error {
	if uid < 1 {
		return errs.ErrorInvalidUid
	}
	idx := uid % s.sessions.length
	if idx >= s.sessions.length {
		return errs.ErrorInvalidSessionSlice
	}
	s.sessions.buckets[idx].lock.Lock()
	defer s.sessions.buckets[idx].lock.Unlock()
	if s.sessions.buckets[idx].sm == nil {
		s.sessions.buckets[idx].sm = make(map[int64]*Session)
	}
	ss, ok := s.sessions.buckets[idx].sm[uid]
	if !ok || ss.Closed {
		return nil
	}
	ss.Close()
	ss.Closed = true
	close(ss.writeChan)
	delete(s.sessions.buckets[idx].sm, uid)
	return nil
}
