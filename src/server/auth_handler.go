package server

import (
	"crypto/md5"
	"fmt"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
)

// http://www.01happy.com/golang-web-get-request-params/
func authHandler(w http.ResponseWriter, r *http.Request, session sessions.Session) {
	pwd := r.PostFormValue("pwd")
	if pwd == currentPassword() {
		salt := fmt.Sprintf("%x", md5.Sum([]byte("auth"+pwd)))
		session.Set("auth", salt)
	} else {
		log.Println("login fail", pwd)
	}
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}
