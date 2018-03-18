package service

import (
	"github.com/DanShu93/martialarts-tracker/storage"
	"github.com/DanShu93/martialarts-tracker/query"
	"reflect"
)

type dummyUUIDGenerator struct {
}

func (g dummyUUIDGenerator) Generate() string {
	return uuidV4Fixture
}

var savedData interface{}
var updatedData interface{}
var deletedData []string
var queriedData query.Query

type dummyRepository struct {
}

func (s dummyRepository) Create(collectionName string, data interface{}) error {
	savedData = data

	return nil
}

func (s dummyRepository) Read(collectionName string, id string, result *interface{}) error {
	if id == missingIDFixture {
		return storage.NotFound
	}

	switch collectionName {
	case reflect.TypeOf(indexedData{}).Name():
		*result = map[string]interface{}{
			"ID":   idFixture,
			"Data": dataValueFixture,
			"NestedData": map[string]interface{}{
				"Data":                    nestedDataValueFixture,
				"DeeplyNestedIndexedData": deeplyNestedIDFixture,
			},
			"NestedIndexedData": nestedIDFixture,
			"SlicedIndexedData": []string{deeplyNestedIDFixture},
		}
	case reflect.TypeOf(nestedIndexedData{}).Name():
		*result = map[string]interface{}{
			"ID":                      nestedIDFixture,
			"Data":                    nestedIndexedDataValueFixture,
			"DeeplyNestedIndexedData": deeplyNestedIDFixture,
		}
	case reflect.TypeOf(deeplyNestedIndexedData{}).Name():
		*result = map[string]interface{}{
			"ID":   deeplyNestedIDFixture,
			"Data": deeplyNestedDataValueFixture,
		}
	default:
		return storage.NotFound
	}

	return nil
}

func (s dummyRepository) Update(collectionName string, id string, data interface{}) error {
	updatedData = data

	return nil
}

func (s dummyRepository) Delete(collectionName string, id string) error {
	switch id {
	case idFixture, nestedIDFixture, deeplyNestedIDFixture:
		deletedData = append(deletedData, id)

		return nil
	}

	return storage.NotFound
}

func (s dummyRepository) ReadAll(collectionName string, query query.Query, result *[]interface{}) error {
	queriedData = query

	var data interface{}

	err := s.Read(collectionName, "", &data)
	if err != nil {
		return err
	}

	*result = []interface{}{data}

	return nil
}
