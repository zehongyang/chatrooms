package models

import (
	"chat_room/db"
	"chat_room/proto/pb"
	"chat_room/utils/errs"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ0123456789"

type DBUser struct {
	dbe *db.DbEngine
}

var GetDBUser = func() func() *DBUser {
	var (
		once sync.Once
		s    *DBUser
	)
	return func() *DBUser {
		once.Do(func() {
			s = &DBUser{
				dbe: db.GetDBEngine(),
			}
		})
		return s
	}
}()

func (s *DBUser) Exists(userName string) (bool, error) {
	return s.dbe.Where("user_name = ?", userName).Exist(new(User))
}

func (s *DBUser) Create(v *User) (int64, error) {
	if v == nil {
		return 0, errs.ErrorNil
	}
	return s.dbe.Insert(v)
}

func (s *DBUser) Update(v *User, cols ...string) (int64, error) {
	if v == nil {
		return 0, errs.ErrorNil
	}
	ss := s.dbe.Where("id = ?", v.Id)
	if len(cols) > 0 {
		ss.Cols(cols...)
	}
	return ss.Update(v)
}

func (s *DBUser) Register(q *pb.UserRegisterQuery, ip string) (*User, error) {
	salt := genSalt(6)
	password := genPassword(q.Password, q.UserName, salt)
	u := &User{
		UserName: q.UserName,
		Token:    genToken(),
		Pwd:      password,
		Salt:     salt,
		NickName: genNickName(),
		Ctm:      time.Now().Unix(),
		Ip:       ip,
	}
	_, err := s.Create(u)
	return u, err
}

func genPassword(password, userName, salt string) string {
	sum := md5.Sum([]byte(fmt.Sprintf("%s%s", password, userName)))
	sum = md5.Sum([]byte(fmt.Sprintf("%s%s", hex.EncodeToString(sum[:]), salt)))
	return hex.EncodeToString(sum[:])
}

func verifyPassword(password string, u *User) bool {
	if u == nil {
		return false
	}
	return genPassword(password, u.UserName, u.Salt) == u.Pwd
}

func genSalt(n int) string {
	var sb strings.Builder
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		idx := rand.Intn(len(letters))
		sb.WriteByte(letters[idx])
	}
	return sb.String()
}

func genToken() string {
	sum := md5.Sum([]byte(genSalt(6)))
	return hex.EncodeToString(sum[:])
}

func genNickName() string {
	return fmt.Sprintf("用户%d%d", time.Now().Unix(), rand.Intn(100))
}
