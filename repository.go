package martialarts

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
)

type TrainingUnitRepository interface {
	Save(trainingUnit TrainingUnit) error
}

type FileTrainingUnitRepository struct {
	DataPath string
}

func (s FileTrainingUnitRepository) Save(trainingUnit TrainingUnit) error {
	series := trainingUnit.Series
	jsonString, err := json.Marshal(trainingUnit)
	if err != nil {
		return err
	}

	seriesDirectoryName := s.DataPath + "/" + series + "/"

	if _, err := os.Stat(seriesDirectoryName); os.IsNotExist(err) {
		err = os.Mkdir(seriesDirectoryName, 0744)
		if err != nil {
			return err
		}
	}

	index, err := getCurrentLessonIndex(seriesDirectoryName)
	if err != nil {
		return err
	}

	filename := index + ".json"
	filePath := seriesDirectoryName + filename

	fmt.Println(filePath)

	return ioutil.WriteFile(filePath, jsonString, 0644)
}

// TODO implement
func getCurrentLessonIndex(seriesDirectoryName string) (string, error) {
	_, err := ioutil.ReadDir(seriesDirectoryName)
	if err != nil {
		return "", err
	}

	return "1", nil
}
