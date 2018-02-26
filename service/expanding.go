package service

import (
	"time"
	"github.com/DanShu93/martialarts-tracker/entity"
)

type ExpandedTrainingUnit struct {
	ID         string `bson:"_id"`
	Start      time.Time
	End        time.Time
	Series     string
	Techniques []entity.Technique
	Methods    []ExpandedMethod
	Exercises  []entity.Exercise
}

type ExpandedMethod struct {
	ID          string `bson:"_id"`
	Kind        string
	Name        string
	Description string
	Covers      []entity.Technique
}

type ExpandingRepository struct {
	TrainingUnitRepository Repository
	TechniqueRepository    Repository
	MethodRepository       Repository
	ExerciseRepository     Repository
}

func (s ExpandingRepository) Save(data interface{}) error {
	return UnsupportedMethod
}

func (s ExpandingRepository) Read(id string, result interface{}) error {
	switch resultPtr := result.(type) {
	case *ExpandedTrainingUnit:
		trainingUnit := entity.TrainingUnit{}
		err := s.TrainingUnitRepository.Read(id, &trainingUnit)
		if err != nil {
			return err
		}

		resultPtr.ID = trainingUnit.ID
		resultPtr.Start = trainingUnit.Start
		resultPtr.End = trainingUnit.End
		resultPtr.Series = trainingUnit.Series

		resultPtr.Techniques = make([]entity.Technique, len(trainingUnit.Techniques))

		for i, techniqueID := range trainingUnit.Techniques {
			technique := entity.Technique{}
			err := s.TechniqueRepository.Read(techniqueID, &technique)
			if err != nil {
				return err
			}

			resultPtr.Techniques[i] = technique
		}

		resultPtr.Exercises = make([]entity.Exercise, len(trainingUnit.Exercises))

		for i, exerciseID := range trainingUnit.Exercises {
			exercise := entity.Exercise{}
			err := s.ExerciseRepository.Read(exerciseID, &exercise)
			if err != nil {
				return err
			}

			resultPtr.Exercises[i] = exercise
		}

		resultPtr.Methods = make([]ExpandedMethod, len(trainingUnit.Methods))

		for i, methodID := range trainingUnit.Methods {
			method := entity.Method{}
			err := s.MethodRepository.Read(methodID, &method)
			if err != nil {
				return err
			}

			expandedMethod := ExpandedMethod{}
			expandedMethod.ID = method.ID
			expandedMethod.Kind = method.Kind
			expandedMethod.Name = method.Name
			expandedMethod.Description = method.Description
			expandedMethod.Covers = make([]entity.Technique, len(method.Covers))

			for i, techniqueID := range method.Covers {
				technique := entity.Technique{}
				err := s.TechniqueRepository.Read(techniqueID, &technique)
				if err != nil {
					return err
				}

				expandedMethod.Covers[i] = technique
			}

			resultPtr.Methods[i] = expandedMethod
		}

	default:
		return UnsupportedEntity
	}

	return nil
}
