package entity

type TrainingUnit struct {
	ID         string `bson:"_id"`
	Series     string
	Techniques []string
	Methods    []string
	Exercises  []string
}

// TODO add time
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
	Covers      []string
}

type Exercise struct {
	ID          string `bson:"_id"`
	Kind        string
	Name        string
	Description string
}
