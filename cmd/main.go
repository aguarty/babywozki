package main

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/jwtauth"

	"go.uber.org/zap"
)

const (
	serviceName string = "babywozki"
)

var (
	version    = "No Version Provided"
	commitHash = "No Git Commit Hash Provided"
)

type application struct {
	logger    *zap.Logger
	cfg       config
	tokenAuth *jwtauth.JWTAuth
	db        *sql.DB
	ctx       context.Context
	tmpls     *appTemplates
}

type appTemplates struct {
	postTmpl      *template.Template
	wozkiListTmpl *template.Template
	errorTmpl     *template.Template
	loginTmpl     *template.Template
	mainTmpl      *template.Template
	appendTmpl    *template.Template
	removeTmpl    *template.Template
}

func main() {

	app := application{}
	if err := app.initConfig("BW"); err != nil {
		panic(err.Error())
	}
	app.tokenAuth = jwtauth.New("HS256", []byte(app.cfg.Secure.Secret), nil)

	app.initLogger(app.cfg.LogLevel, app.cfg.LogType)
	defer app.logger.Sync()

	if err := app.initDB(); err != nil {
		app.logger.Fatal("DATABASE", zap.Error(err))
	}

	if err := app.initTemplates(); err != nil {
		app.logger.Fatal("TEMPLATES", zap.Error(err))
	}
	// create context
	ctx, cancelWork := context.WithCancel(context.Background())
	app.ctx = ctx

	handler, err := app.createHTTPHandler()
	if err != nil {
		app.logger.Fatal("HANDLER", zap.Error(err))
	}

	// create server
	server := &http.Server{
		Addr:         app.cfg.Server.Host + ":" + app.cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  app.cfg.Web.ReadTimeout,
		WriteTimeout: app.cfg.Web.WriteTimeout,
		IdleTimeout:  app.cfg.Web.ReadTimeout,
	}

	// start server
	serverErrors := make(chan error, 1)
	go func() {
		app.logger.Info("SERVER", zap.String("Listen and serve host", server.Addr))
		serverErrors <- server.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	// graceful shutdown
	select {
	case err := <-serverErrors:
		app.logger.Fatal("SERVER", zap.String("msg", "Can`t start server"), zap.Error(err))
	case <-osSignals:
		app.logger.Info("SERVER", zap.String("msg", "Server shutdown"))
		ctx, cancelServ := context.WithTimeout(context.Background(), app.cfg.Web.ShutdownTimeout)
		if err := server.Shutdown(ctx); err != nil {
			app.logger.Error("SERVER", zap.String("msg", "Graceful shutdown did not complete"), zap.Error(err))
			if err := server.Close(); err != nil {
				app.logger.Fatal("SERVER", zap.String("msg", "Could not stop http server"), zap.Error(err))
			}
		}
		cancelServ()
		cancelWork()
	}
	app.logger.Info("SERVER", zap.String("msg", "Service stopped"))
}
