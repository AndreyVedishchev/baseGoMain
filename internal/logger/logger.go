package logger

import "go.uber.org/zap"

// Init создаёт development-логгер zap и настраивает вывод в файл app.log
func Init() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"app.log"}
	cfg.ErrorOutputPaths = []string{"app.log"}

	return cfg.Build()
}
