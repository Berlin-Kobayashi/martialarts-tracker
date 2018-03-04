package entity

import "time"

type TrainingUnit struct {
	ID         string `bson:"_id"`
	Start      time.Time
	End        time.Time
	Series     string
	Techniques []Technique
	Methods    []Method
	Exercises  []Exercise
}

type Technique struct {
	ID          string `bson:"_id"`
	Kind        string
	Name        string
	Description string
}

type Method struct {
	ID          string `bson:"_id"`
	Kind        string
	Name        string
	Description string
	Covers      []Technique
}

type Exercise struct {
	ID          string `bson:"_id"`
	Kind        string
	Name        string
	Description string
}
