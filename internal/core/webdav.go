package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"smartcloud/internal/config"

	"golang.org/x/net/webdav"
)

type WebDav struct {
	permission *Permission
	handler    *webdav.Handler
}

func (web *WebDav) dir(w http.ResponseWriter, req *http.Request) bool {
	ctx := context.Background()

	// web.handler.FileSystem.
	f, err := web.handler.FileSystem.OpenFile(ctx, req.URL.Path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer f.Close()
	if fi, _ := f.Stat(); fi != nil && !fi.IsDir() {
		return false
	}
	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Print(w, "Error reading directory", http.StatusInternalServerError)
		return false
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", name, name)
	}
	fmt.Fprintf(w, "</pre>\n")
	return true
}

func NewWebDav(conf config.Permission, fs *webdav.Handler) *WebDav {

	return &WebDav{
		permission: NewPermission(conf),
		handler:    fs,
	}
}

func (web *WebDav) DoService(w http.ResponseWriter, r *http.Request) {

	// user := r.Header.Get("user")

	user := "user"
	if !web.permission.Enforce(user, r.URL.Path, r.Method) {
		log.Printf("permission verify failed %s,%s,%s", user, r.URL.Path, r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// if r.Method == "GET" && web.dir(w, r) {
	// 	return
	// }

	// switch r.Method {
	// case http.MethodOptions:
	// 	handleOptions(w, r)
	// case http.MethodGet:
	// 	handleGet(w, r)
	// case http.MethodPut:
	// 	web.handlePut(w, r)
	// case http.MethodDelete:
	// 	handleDelete(w, r)
	// default:
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// }
}
