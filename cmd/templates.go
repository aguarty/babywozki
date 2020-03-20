package main

import (
	"html/template"
	"path"
)

func (a *application) initTemplates() (err error) {
	a.tmpls = &appTemplates{}

	a.tmpls.mainTmpl, err = template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "index.html"))
	if err != nil {
		return
	}
	a.tmpls.postTmpl, err = template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "post.html"))
	if err != nil {
		return
	}
	a.tmpls.errorTmpl, err = template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "error.html"))
	if err != nil {
		return
	}
	a.tmpls.wozkiListTmpl, err = template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "wozkilist.html"))
	if err != nil {
		return
	}
	a.tmpls.loginTmpl, err = template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "login.html"))
	if err != nil {
		return
	}
	a.tmpls.appendTmpl, err = template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "append.html"))
	if err != nil {
		return
	}
	a.tmpls.removeTmpl, err = template.ParseFiles(path.Join("templates", "layout.html"), path.Join("templates", "remove.html"))
	if err != nil {
		return
	}
	return
}
