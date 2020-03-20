package main

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func (a *application) initDB() error {

	src := "host=" + a.cfg.DB.Host + " port=" + a.cfg.DB.Port + " sslmode=disable dbname=" + a.cfg.DB.Database + " user=" + a.cfg.DB.User + " password=" + a.cfg.DB.Password

	db, err := sql.Open("postgres", src)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	//db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(time.Second * 10)
	a.db = db
	return nil
}
