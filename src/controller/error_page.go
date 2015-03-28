package controller

import (
	"io"
	"net/http"
)

func ErrorPage(w http.ResponseWriter, status int, errmsg string) {
	w.WriteHeader(status)
	io.WriteString(w, errmsg)
}
