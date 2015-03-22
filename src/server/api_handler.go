package server

import (
	"log"
	"net/http"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())

	if r.URL.String() == "/api/rebuild" {
		rebuildAll()
	}
}

func rebuildAll() {

	blogBuilder.RebuildAll()
}
