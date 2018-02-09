package martialarts

import (
	"encoding/json"
	"io/ioutil"
)

type TrainingUnitRepository interface {
	Save(trainingUnit TrainingUnit) error
}

type FileTrainingUnitRepository struct {
}

func (s FileTrainingUnitRepository) Save(trainingUnit TrainingUnit) error {
	jsonString, err := json.Marshal(trainingUnit)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("/go/src/github.com/DanShu93/martialarts-tracker/data/training1.json", jsonString, 0644)

}
