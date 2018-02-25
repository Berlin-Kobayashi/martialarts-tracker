package service

import (
	"net/http"
)

type LogService struct {
}

func (s LogService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.post(rw, r)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

//TODO implement
func (s LogService) post(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("logged"))
}
