package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Server struct {
		Host string `envconfig:"SERVER_HOST" default:""`
		Port string `envconfig:"SERVER_PORT" default:"8888"`
	}
	DB struct {
		Host     string        `default:"localhost"`
		Port     string        `default:"5432"`
		Database string        `default:""`
		User     string        `default:""`
		Password string        `default:""`
		MAX      int           `default:"3"`
		TTL      time.Duration `default:"10s"`
	}
	Web struct {
		MaxHeaderBytes  int           `envconfig:"WEB_MAX_HEADER_BYTES" default:"16777216"` // 16Mb
		ReadTimeout     time.Duration `envconfig:"WEB_READ_TIMEOUT" default:"10s"`
		WriteTimeout    time.Duration `envconfig:"WEB_WRITE_TIMEOUT" default:"10s"`
		ShutdownTimeout time.Duration `envconfig:"WEB_SHUTDOWN_TIMEOUT" default:"10s"`
	}
	Secure struct {
		Secret string `default:""`
		Salt   string `default:""`
	}
	// log level
	LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
	// формат вывода лога: json или console
	LogType string `envconfig:"LOG_TYPE" default:"json"`
}

func (a *application) initConfig(app string) error {
	var conf config
	if err := envconfig.Process(app, &conf); err != nil {
		envconfig.Usage(app, &conf)
		return err
	}
	a.cfg = conf
	return nil
}
