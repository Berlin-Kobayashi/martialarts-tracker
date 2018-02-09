package main

import (
	"net/http"
	"github.com/DanShu93/martialarts-tracker"
)

func main() {
	repository := martialarts.FileTrainingUnitRepository{}

	http.Handle("/training-unit", martialarts.TrainingUnitService{Repository: repository})
	http.ListenAndServe(":80", nil)
}
