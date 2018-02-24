package repository

import "github.com/DanShu93/martialarts-tracker/entity"

var RecordedTrainingUnit entity.TrainingUnit

var TrainingUnitFixture = entity.TrainingUnit{
	Series: "JKD I",
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

type DummyTrainingUnitRepository struct {
}

func (s DummyTrainingUnitRepository) Save(trainingUnit entity.TrainingUnit) error {
	RecordedTrainingUnit = trainingUnit

	return nil
}

func (s DummyTrainingUnitRepository) Read(trainingUnitIndex string) (entity.TrainingUnit, error) {
	return TrainingUnitFixture, nil
}
