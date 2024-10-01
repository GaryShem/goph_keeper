package logging

import (
	"sync"

	"go.uber.org/zap"
)

var log *zap.SugaredLogger = nil
var logOnce sync.Once

func Log() *zap.SugaredLogger {
	level := "info"
	logOnce.Do(func() {
		lvl, err := zap.ParseAtomicLevel(level)
		if err != nil {
			panic(err)
		}
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = lvl
		zl, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		log = zl.Sugar()
	})

	return log
}
