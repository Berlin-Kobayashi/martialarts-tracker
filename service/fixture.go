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
	Techniques: []string{
		"bc5ac88f-3d3f-4a1a-83b2-92f847eb6ae6",
	},
	Methods: []string{
		"2969227d-6eb4-4cdd-96f5-7ca6c97d4df8",
	},
	Exercises: []string{
		"4da2ab63-d83d-4b14-8cf0-e0d8eb815ce8",
	},
}
