package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (a *application) initLogger(loglevel, logtype string) {

	level := zap.NewAtomicLevel()
	switch loglevel {
	case "debug":
		level.SetLevel(zapcore.DebugLevel)
	case "info":
		level.SetLevel(zapcore.InfoLevel)
	case "warn":
		level.SetLevel(zapcore.WarnLevel)
	case "error":
		level.SetLevel(zapcore.ErrorLevel)
	default:
		level.SetLevel(zapcore.InfoLevel)
	}

	// Можно сделать вывод в файл
	//===========================
	// var OutPath []string
	// if *logfile != "" {
	// 	OutPath = []string{"stdout", logdir + *logfile}
	// } else {
	// 	OutPath = []string{"stdout"}
	// }

	a.logger, _ = zap.Config{
		Level:    level,
		Encoding: logtype,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:   "Time",
			LevelKey:  "Level",
			NameKey:   "Logger",
			CallerKey: "Source",
			//StacktraceKey: "St",
			MessageKey:     "TYPE",
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
}
