package martialarts

type TrainingUnit struct {
	Series     string
	Techniques []Technique
	Methods    []Method
	Exercises  []Exercise
}

type Technique struct {
	Kind        string
	Name        string
	Description string
}

type Method struct {
	Kind        string
	Name        string
	Description string
	Covers      []Technique
}

type Exercise struct {
	Kind        string
	Name        string
	Description string
}
