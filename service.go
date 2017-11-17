package main

import (
	"net/http"
)

type MartialartsTrackerService struct {
}

func (s MartialartsTrackerService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		rw.Write([]byte("JKD"))
	}
}

func main() {
	http.ListenAndServe(":80", MartialartsTrackerService{})
}
