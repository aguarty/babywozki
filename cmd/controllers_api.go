package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"

	"go.uber.org/zap"
)

// appendItem -
func (a *application) appendItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		file, handler, err := r.FormFile("imgfile")
		if err != nil {
			a.logger.Error("APPEND", zap.Error(err))
			sendResponse(a.logger, w, http.StatusInternalServerError, ApiResp{Err: err.Error()})
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			a.logger.Error("STORAGE", zap.Error(err))
		}

		brand := r.FormValue("brand")
		if brand == "" {
			sendResponse(a.logger, w, http.StatusOK, ApiResp{Err: "Invalid input"})
			return
		}
		model := r.FormValue("model")
		if model == "" {
			sendResponse(a.logger, w, http.StatusOK, ApiResp{Err: "Invalid input"})
			return
		}
		price, err := strconv.Atoi(r.FormValue("price"))
		if err != nil {
			sendResponse(a.logger, w, http.StatusOK, ApiResp{Err: "Invalid input"})
			return
		}
		description := r.FormValue("description")
		if description == "" {
			sendResponse(a.logger, w, http.StatusOK, ApiResp{Err: "Invalid input"})
			return
		}

		ext := string(filepath.Ext(handler.Filename)[1:])
		_, err = a.db.ExecContext(a.ctx, `INSERT INTO babywozki.wozki(brand, model ,description ,price ,img_name ,img ,img_ext ) 
		VALUES ($1, $2, $3, $4, $5, $6::bytea, $7)`, brand, model, description, price, handler.Filename, base64.StdEncoding.EncodeToString(fileBytes), ext)
		if err != nil {
			a.logger.Error("ERROR", zap.Error(err))
			sendResponse(a.logger, w, http.StatusInternalServerError, ApiResp{Err: err.Error()})
			return
		}

		sendResponse(a.logger, w, http.StatusOK, ApiResp{Data: "OK"})
	}
}

// removeItem -
func (a *application) removeItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		strId := chi.URLParam(r, "itemID")
		if strId == "" {
			sendResponse(a.logger, w, http.StatusOK, ApiResp{Err: "Invalid input"})
			return
		}

		id, err := strconv.Atoi(strId)
		if err != nil {
			sendResponse(a.logger, w, http.StatusOK, ApiResp{Err: "Invalid input"})
			return
		}

		_, err = a.db.ExecContext(a.ctx, `delete from babywozki.wozki where id=$1;`, id)
		if err != nil {
			a.logger.Error("ERROR", zap.Error(err))
			sendResponse(a.logger, w, http.StatusInternalServerError, ApiResp{Err: err.Error()})
			return
		}
		sendResponse(a.logger, w, http.StatusOK, ApiResp{Data: "OK"})
	}
}
