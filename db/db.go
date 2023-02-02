package db

import (
	"chat_room/utils/config"
	"chat_room/utils/logger"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"sync"
	"xorm.io/xorm"
)

func init() {
	GetDBEngine()
}

type DbEngine struct {
	*xorm.Engine
}

type DbConfig struct {
	Db struct {
		Driver   string
		Uri      string
		MaxNums  int
		IdleNums int
	}
}

var GetDBEngine = func() func() *DbEngine {
	var (
		once sync.Once
		dbe  *DbEngine
		cfg  DbConfig
	)
	return func() *DbEngine {
		once.Do(func() {
			err := config.Load(&cfg)
			if err != nil {
				logger.Fatal("GetDBEngine", zap.Error(err))
			}
			engine, err := xorm.NewEngine(cfg.Db.Driver, cfg.Db.Uri)
			if err != nil {
				logger.Fatal("GetDBEngine", zap.Error(err), zap.Any("cfg", cfg))
			}
			engine.SetMaxOpenConns(cfg.Db.MaxNums)
			engine.SetMaxIdleConns(cfg.Db.IdleNums)
			err = engine.Ping()
			if err != nil {
				logger.Fatal("GetDBEngine", zap.Error(err), zap.Any("cfg", cfg))
			}
			dbe = &DbEngine{engine}
		})
		return dbe
	}
}()
