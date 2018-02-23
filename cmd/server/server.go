package main

import (
	"net/http"
	"github.com/DanShu93/martialarts-tracker/repository"
	"github.com/DanShu93/martialarts-tracker/service"
)

func main() {
	repo := repository.FileTrainingUnitRepository{DataPath: "/go/src/github.com/DanShu93/martialarts-tracker/data/trainingunits"}

	http.Handle("/training-unit/", service.TrainingUnitService{Repository: repo})
	http.Handle("/training-unit", service.TrainingUnitService{Repository: repo})

	http.ListenAndServe(":80", nil)
}
