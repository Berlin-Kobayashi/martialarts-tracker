package service

import (
	"github.com/DanShu93/martialarts-tracker/entity"
	"time"
)

var uuidV4Fixture = "b5e57615-0f40-404e-bbe0-6ae81fe8080a"

var trainingUnitFixture = entity.TrainingUnit{
	Series: "Coding",
	Start:  time.Date(2018, 2, 25, 17, 0, 0, 0, time.UTC),
	End:    time.Date(2018, 2, 25, 17, 0, 0, 0, time.UTC),
	ID:     "b5e57615-0f40-404e-bbe0-6ae81fe8080a",
	Techniques: []entity.Technique{
		{
			ID:          "bc5ac88f-3d3f-4a1a-83b2-92f847eb6ae6",
			Kind:        "Language",
			Name:        "Go",
			Description: "compiled, concurrent, imperative, structured",
		},
	},
	Methods: []entity.Method{
		{
			ID:          "2969227d-6eb4-4cdd-96f5-7ca6c97d4df8",
			Kind:        "Project",
			Name:        "martialarts-tracker",
			Description: "Tracking service for martial arts training.",
			Covers: []entity.Technique{
				{
					ID:          "bc5ac88f-3d3f-4a1a-83b2-92f847eb6ae6",
					Kind:        "Language",
					Name:        "Go",
					Description: "compiled, concurrent, imperative, structured",
				},
			},
		},
	},
	Exercises: []entity.Exercise{
		{
			ID:          "4da2ab63-d83d-4b14-8cf0-e0d8eb815ce8",
			Kind:        "Software package",
			Name:        "Filesystem database",
			Description: "Persist and structure data in a filesystem.",
		},
	},
}
