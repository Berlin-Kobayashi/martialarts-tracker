package repository

import "github.com/DanShu93/martialarts-tracker/entity"

var RecordedTrainingUnit entity.TrainingUnit

var pakSaoFixture = entity.Technique{
	Kind:        "counter",
	Name:        "pak sao",
	Description: "Means slapping hand\nCounter a jab by slapping the elbow of the opponent into his body, destroying his structure\nAt the same time perform a jab",
}

var TrainingUnitFixture = entity.TrainingUnit{
	Series: "JKD I",
	Techniques: []entity.Technique{
		pakSaoFixture,
	},
	Methods: []entity.Method{
		{
			Kind:        "counter",
			Name:        "Pak Sao drill",
			Description: "",
			Covers:      []entity.Technique{pakSaoFixture},
		},
	},
	Exercises: []entity.Exercise{
		{
			Kind:        "Sparring",
			Name:        "Lead hand sparring",
			Description: "Sparing with lead hand punches only",
		},
	},
}

type DummyTrainingUnitRepository struct {
}

func (s DummyTrainingUnitRepository) Save(trainingUnit entity.TrainingUnit) (string, error) {
	RecordedTrainingUnit = trainingUnit

	return "1", nil
}

func (s DummyTrainingUnitRepository) Read(trainingSeriesName, trainingUnitIndex string) (entity.TrainingUnit, error) {
	return TrainingUnitFixture, nil
}
