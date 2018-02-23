package martialarts

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"regexp"
	"errors"
)

var NotFound = errors.New("training unit could not be found")
var Invalid = errors.New("training unit is invalid")

type TrainingUnitRepository interface {
	Save(trainingUnit TrainingUnit) (string, error)
	Read(trainingSeriesName, trainingUnitIndex string) (TrainingUnit, error)
}

type FileTrainingUnitRepository struct {
	DataPath string
}

func (r FileTrainingUnitRepository) Save(trainingUnit TrainingUnit) (string, error) {
	series := trainingUnit.Series
	jsonString, err := json.Marshal(trainingUnit)
	if err != nil {
		return "", err
	}

	seriesDirectoryName := r.DataPath + "/" + series + "/"

	if _, err := os.Stat(seriesDirectoryName); os.IsNotExist(err) {
		err = os.Mkdir(seriesDirectoryName, 0744)
		if err != nil {
			return "", err
		}
	}

	index, err := getCurrentLessonIndex(seriesDirectoryName)
	if err != nil {
		return "", err
	}

	filePath := r.getFilePath(series, index)

	err = ioutil.WriteFile(filePath, jsonString, 0644)
	if err != nil {
		return "", err
	}

	return index, nil
}

func (r FileTrainingUnitRepository) getFilePath(trainingSeriesName, trainingUnitIndex string) string {
	return r.DataPath + "/" + trainingSeriesName + "/" + trainingUnitIndex + ".json"
}

func getCurrentLessonIndex(seriesDirectoryName string) (string, error) {
	files, err := ioutil.ReadDir(seriesDirectoryName)
	if err != nil {
		return "", err
	}

	maxIndex := 0

	fileNameRegex := regexp.MustCompile("^(\\d+)\\..*$")
	for _, f := range files {
		indexString := fileNameRegex.ReplaceAll([]byte(f.Name()), []byte("$1"))
		index, err := strconv.Atoi(string(indexString))
		if err != nil {
			return "", err
		}

		if index > maxIndex {
			maxIndex = index
		}
	}

	maxIndex++

	return strconv.Itoa(maxIndex), nil
}

func (r FileTrainingUnitRepository) Read(trainingSeriesName, trainingUnitIndex string) (TrainingUnit, error) {
	filePath := r.getFilePath(trainingSeriesName, trainingUnitIndex)

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return TrainingUnit{}, NotFound
	}

	trainingUnit := TrainingUnit{}
	err = json.Unmarshal(contents, &trainingUnit)
	if err != nil {
		return TrainingUnit{}, Invalid
	}

	return trainingUnit, nil
}
