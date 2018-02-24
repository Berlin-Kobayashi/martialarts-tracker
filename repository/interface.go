package repository

import "github.com/DanShu93/martialarts-tracker/entity"

type TrainingUnitRepository interface {
	Save(trainingUnit entity.TrainingUnit) error
	Read(trainingUnitIndex string) (entity.TrainingUnit, error)
}
