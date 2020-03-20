package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//createHTTPHandler create handler
func (app *application) createHTTPHandler() (http.Handler, error) {

	// serve static
	fs := http.FileServer(http.Dir("./static"))
	mux := chi.NewMux()

	mux.Use(middleware.Recoverer)
	mux.Use(app.logging())
	mux.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cache-Control", "Pragma"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	// serve static
	mux.Handle("/static/*", http.StripPrefix("/static/", fs))
	// public pages
	mux.Route("/", func(main chi.Router) {
		main.Use(middleware.StripSlashes)
		main.Get("/", app.mainPage())
		main.Get("/login", app.loginPage())
		main.Get("/wozki", app.wozkiPage())
		main.Get("/wozki/{brand}", nil)
		main.Get("/wozki/{brand}/{id}", nil)
		main.Get("/{anytext}", app.errorHandlerCode(404))
	})

	// admin pages
	mux.Route("/admin", func(admin chi.Router) {
		admin.Use(Verifier(app.tokenAuth, app.cfg.Secure.Salt))
		admin.Use(app.Authenticator)
		admin.Get("/append", app.appendPage())
		admin.Get("/remove", app.removePage())
		admin.Get("/edit", nil)

		admin.Handle("/upload", app.uploadImg())
		admin.Handle("/metrics", promhttp.Handler())

	})

	// API
	mux.Route("/api", func(api chi.Router) {
		api.Post("/login", app.loginHandler())
		api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			sendResponse(app.logger, w, http.StatusOK, HealthResponse{
				Code:       http.StatusOK,
				Version:    version,
				CommitHash: commitHash,
			})
		})
		api.Route("/v1", func(v1 chi.Router) {
			v1.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8;"))
			v1.Use(Verifier(app.tokenAuth, app.cfg.Secure.Salt))
			v1.Use(app.Authenticator)
			v1.Route("/wozki", func(wozki chi.Router) {
				wozki.Post("/append", app.appendItem())
				wozki.Get("/remove/{itemID}", app.removeItem())
			})
		})
	})

	return mux, nil
}

//logging - middleware for logging
func (a *application) logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			// --- metrics -----------------------------------------------------
			//path := strings.Trim(r.URL.Path, "/")
			requestsCounter.WithLabelValues(r.URL.Path).Inc()
			defer func() {
				latency.WithLabelValues(r.URL.Path).Observe(time.Since(start).Seconds())
				responsesCounter.WithLabelValues(r.URL.Path, fmt.Sprintf("%d", ww.Status())).Inc()
			}()
			// -----------------------------------------------------------------
			next.ServeHTTP(ww, r)
			fields := []zapcore.Field{
				zap.Int("code", ww.Status()),
				zap.String("Method", r.Method),
				zap.String("url", r.URL.Path),
				zap.String("addr", r.RemoteAddr),
				zap.Int64("req_lenght", r.ContentLength),
				zap.String("latency", time.Since(start).String()),
				zap.Int("bytes", ww.BytesWritten()),
			}
			defer func() {
				switch {
				case ww.Status() >= 500:
					fields = append(fields, zap.String("msg", "Server error"))
					a.logger.Error("WEB", fields...)
				case ww.Status() >= 400:
					fields = append(fields, zap.String("msg", "Client error"))
					a.logger.Warn("WEB", fields...)
				case ww.Status() >= 300:
					fields = append(fields, zap.String("msg", "Redirection"))
					a.logger.Info("WEB", fields...)
				default:
					fields = append(fields, zap.String("msg", "Success"))
					a.logger.Debug("WEB", fields...)
				}
			}()
		})
	}
}
