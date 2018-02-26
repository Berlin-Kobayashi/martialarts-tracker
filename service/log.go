package service

import (
	"time"
	"github.com/DanShu93/martialarts-tracker/entity"
)

type Log struct {
	ID         string `bson:"_id"`
	Start      time.Time
	End        time.Time
	Series     string
	Techniques []entity.Technique
	Methods    []entity.Method
	Exercises  []entity.Exercise
}

type LogRepository struct {
}

func (s LogRepository) Save(data interface{}) error {
	return nil
}

func (s LogRepository) Read(id string, result interface{}) error {
	return nil
}
