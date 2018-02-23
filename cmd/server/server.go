package main

import (
	"net/http"
	"github.com/DanShu93/martialarts-tracker"
)

func main() {
	repository := martialarts.FileTrainingUnitRepository{DataPath: "/go/src/github.com/DanShu93/martialarts-tracker/data/trainingunits"}

	http.Handle("/training-unit/", martialarts.TrainingUnitService{Repository: repository})
	http.Handle("/training-unit", martialarts.TrainingUnitService{Repository: repository})

	http.ListenAndServe(":80", nil)
}
