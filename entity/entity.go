package entity

type TrainingUnit struct {
	ID         string
	Start      string
	End        string
	Series     string
	Techniques []Technique
	Methods    []Method
	Exercises  []Exercise
}

type Technique struct {
	ID          string
	Kind        string
	Name        string
	Description string
}

type Method struct {
	ID          string
	Kind        string
	Name        string
	Description string
	Covers      []Technique
}

type Exercise struct {
	ID          string
	Kind        string
	Name        string
	Description string
}
