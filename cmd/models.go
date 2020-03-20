package main

import "html/template"

// HealthResponse -
type HealthResponse struct {
	Code       int    `json:"code,omitempty"`
	Date       string `json:"date,omitempty"`
	Error      string `json:"error,omitempty"`
	Version    string `json:"version,omitempty"`
	CommitHash string `json:"commit,omitempty"`
	Info       string `json:"information,omitempty"`
}

// ApiResp -
type ApiResp struct {
	Err   string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Count int         `json:"count,omitempty"`
}

type Posts []Post

type Post struct {
	Title string
	Body  template.HTML
	Img   string
}

type Wozki []Wozka

type Wozka struct {
	Id          int64
	Brand       string
	Description string
	Model       string
	Img         string
	Ext         string
	Price       int64
}

type image struct {
	id      int64
	imgName string
	img     []byte
	ext     string
}

type DbWozka struct {
	Id          int64
	Brand       string
	Model       string
	Description string
	Price       int
	ImgName     string
	Img         []byte
	ImgExt      string
}
