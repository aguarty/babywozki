package main

import (
	"database/sql"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"go.uber.org/zap"
)

// wozkiPage -
func (a *application) uploadImg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		dir := "./static/images/site/"
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			a.logger.Error("ERROR", zap.Error(err))
			sendResponse(a.logger, w, http.StatusInternalServerError, ApiResp{Err: err.Error()})
			return
		}

		for _, file := range files {
			f, err := ioutil.ReadFile(dir + file.Name())
			if err != nil {
				a.logger.Error("ERROR", zap.Error(err))
				sendResponse(a.logger, w, http.StatusInternalServerError, ApiResp{Err: err.Error()})
				return
			}
			ext := string(filepath.Ext(file.Name())[1:])
			_, err = a.db.ExecContext(a.ctx, "INSERT INTO babywozki.images VALUES ($1, $2, $3::bytea, $4)", 123, file.Name(), base64.StdEncoding.EncodeToString(f), ext)
			if err != nil {
				a.logger.Error("ERROR", zap.Error(err))
				sendResponse(a.logger, w, http.StatusInternalServerError, ApiResp{Err: err.Error()})
				return
			}
		}
	}
}

// loginPage -
func (a *application) loginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := a.tmpls.loginTmpl.ExecuteTemplate(w, "layout", nil); err != nil {
			a.logger.Error("LOGIN_PAGE", zap.Error(err))
			a.errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}
}

// mainPage -
func (a *application) mainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := a.tmpls.mainTmpl.ExecuteTemplate(w, "layout", nil); err != nil {
			a.logger.Error("MAIN_PAGE", zap.Error(err))
			a.errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}
}

// appendPage -
func (a *application) appendPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := a.tmpls.appendTmpl.ExecuteTemplate(w, "layout", nil); err != nil {
			a.logger.Error("APPEND_PAGE", zap.Error(err))
			a.errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}
}

// removePage -
func (a *application) removePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		query := "select id, brand, model, price from babywozki.wozki;"
		rows, err := a.db.QueryContext(a.ctx, query)
		if err != nil {
			a.logger.Error("ERROR", zap.Error(err))
			sendResponse(a.logger, w, http.StatusInternalServerError, ApiResp{Err: err.Error()})
			return
		}
		data := make([]DbWozka, 0)
		for rows.Next() {
			tmp := DbWozka{}
			if err := rows.Scan(
				&tmp.Id,
				&tmp.Brand,
				&tmp.Model,
				&tmp.Price,
			); err != nil {
				rows.Close()
				a.logger.Error("ERROR", zap.Error(err))
				sendResponse(a.logger, w, http.StatusInternalServerError, ApiResp{Err: err.Error()})
				return
			}
			data = append(data, tmp)
		}

		if err := a.tmpls.removeTmpl.ExecuteTemplate(w, "layout", data); err != nil {
			a.logger.Error("REMOVE_PAGE", zap.Error(err))
			a.errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}
}

// wozkiPage -
func (a *application) wozkiPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		wozki, status, err := loadWozki(a.db)
		if err != nil {
			a.logger.Error("WOZKI_PAGE", zap.Error(err))
			a.errorHandler(w, r, status)
			return
		}

		if err := a.tmpls.wozkiListTmpl.ExecuteTemplate(w, "layout", wozki); err != nil {
			a.logger.Error("WOZKI_PAGE", zap.Error(err))
			a.errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}
}

func (a *application) errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if err := a.tmpls.errorTmpl.ExecuteTemplate(w, "layout", map[string]interface{}{"Error": http.StatusText(status), "Status": status}); err != nil {
		a.logger.Error("ERROR", zap.Error(err))
		sendResponse(a.logger, w, http.StatusInternalServerError, ApiResp{Err: err.Error()})
		return
	}
}

// errorHandlerCode -
func (a *application) errorHandlerCode(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.errorHandler(w, r, status)
	}
}

// loadWozki -
func loadWozki(db *sql.DB) (Wozki, int, error) {

	img := image{}
	err := db.QueryRow("select * from babywozki.images where id=$1", 123).Scan(&img.id, &img.imgName, &img.img, &img.ext)
	if err != nil {
		return []Wozka{}, 500, err
	}
	var wozki Wozki
	for i := 0; i < 5; i++ {
		wozki = append(wozki, Wozka{
			Name:        "Adamex ASD324 002",
			Description: "Это просто рыба. Это просто рыба. Это просто рыба. Это просто рыба. Это просто рыба. Это просто рыба. Это просто рыба.",
			Img:         string(img.img),
			Price:       int64(i),
			Ext:         img.ext,
		})
	}

	return wozki, 200, nil
}
