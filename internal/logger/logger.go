package logger

import "go.uber.org/zap"

// Log — глобальный логгер приложения, инициализируется вызовом Init()
var Log *zap.Logger

// Init создаёт development-логгер zap и настраивает вывод в файл app.log
func Init() error {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"app.log"}
	cfg.ErrorOutputPaths = []string{"app.log"}

	l, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = l
	return nil
}
