package service

import "reflect"

type IndexingRepository struct {
	Generator  UUIDGenerator
	Repository Repository
}

func (s IndexingRepository) Save(data interface{}) error {
	id := s.Generator.Generate()
	reflect.ValueOf(data).Elem().FieldByName("ID").SetString(id)

	return s.Repository.Save(data)
}

func (s IndexingRepository) Read(id string, result interface{}) error {
	return s.Repository.Read(id, result)
}
