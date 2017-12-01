package main

import (
	"net/http"
)

type FavIconService struct {
}

func (s FavIconService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "image/x-icon")
	http.ServeFile(rw, r, "/go/src/github.com/DanShu93/martialarts-tracker/favicon.ico")
}

type MainService struct {
}

func (s MainService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/html")
	rw.Write([]byte("<link rel=\"shortcut icon\" href=\"http://localhost:8888/favicon.ico\" type=\"image/x-icon\">"))
}


func main() {
	http.Handle("/favicon.ico", FavIconService{})
	http.Handle("/index.html", MainService{})
	http.ListenAndServe(":80", nil)
}
