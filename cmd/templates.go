package main

import (
	"html/template"
	"path"
)

func (a *application) initTemplates() (err error) {
	a.tmpls = &appTemplates{}
	a.tmpls.mainTmpl, err = template.ParseFiles(
		path.Join("templates", "layout.html"),
		path.Join("templates", "index.html"),
		path.Join("templates", "not-admin.tmpl"),
	)
	if err != nil {
		return
	}
	a.tmpls.postTmpl, err = template.ParseFiles(
		path.Join("templates", "layout.html"),
		path.Join("templates", "post.html"),
		path.Join("templates", "not-admin.tmpl"),
	)
	if err != nil {
		return
	}
	a.tmpls.errorTmpl, err = template.ParseFiles(
		path.Join("templates", "layout.html"),
		path.Join("templates", "error.html"),
		path.Join("templates", "not-admin.tmpl"),
	)
	if err != nil {
		return
	}
	a.tmpls.wozkiListTmpl, err = template.ParseFiles(
		path.Join("templates", "layout.html"),
		path.Join("templates", "wozkilist.html"),
		path.Join("templates", "not-admin.tmpl"),
	)
	if err != nil {
		return
	}
	a.tmpls.loginTmpl, err = template.ParseFiles(
		path.Join("templates", "layout.html"),
		path.Join("templates", "login.html"),
		path.Join("templates", "not-admin.tmpl"),
	)
	if err != nil {
		return
	}
	a.tmpls.appendTmpl, err = template.ParseFiles(
		path.Join("templates", "layout.html"),
		path.Join("templates", "append.html"),
		path.Join("templates", "admin-menu.tmpl"),
	)
	if err != nil {
		return
	}
	a.tmpls.removeTmpl, err = template.ParseFiles(
		path.Join("templates", "layout.html"),
		path.Join("templates", "remove.html"),
		path.Join("templates", "admin-menu.tmpl"),
	)
	if err != nil {
		return
	}
	return
}
