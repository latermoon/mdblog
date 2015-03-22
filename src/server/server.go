package server

import (
	"builder"
	"log"
	"net/http"
	"path/filepath"
)

var Workspace string // website home directory
var blogBuilder *builder.BlogBuilder

func ListenAndServe(host string, workspace string) {
	Workspace = workspace

	// builder
	var err error
	blogBuilder, err = builder.NewBlogBuilder(Workspace)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)

	// http server
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/img/", imageHandler)
	go http.ListenAndServe(host, nil)

	log.Println("watching:", workspace)
	go watch(filepath.Join(workspace, "article"), filepath.Join(workspace, "template"))

	<-done
}
