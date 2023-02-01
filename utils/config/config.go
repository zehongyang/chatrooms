package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var gData []byte
var GLog *zap.Logger

type LogConfig struct {
	Zap zap.Config
}

func init() {
	data, err := os.ReadFile("./application.yml")
	if err != nil {
		panic(err)
	}
	gData = data
	err = initLogger()
	if err != nil {
		panic(err)
	}
}

func initLogger() error {
	tlog, err := GetLogConfig().Zap.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	GLog = tlog
	return nil
}

var GetLogConfig = func() func() LogConfig {
	var (
		once sync.Once
		cfg  LogConfig
	)
	return func() LogConfig {
		once.Do(func() {
			err := Load(&cfg)
			if err != nil {
				panic(err)
			}
		})
		return cfg
	}
}()

func Load(v interface{}) error {
	if len(gData) < 1 {
		return nil
	}
	return yaml.Unmarshal(gData, v)
}
