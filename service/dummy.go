package service

import (
	"github.com/DanShu93/martialarts-tracker/storage"
)

type dummyUUIDGenerator struct {
}

func (g dummyUUIDGenerator) Generate() string {
	return uuidV4Fixture
}

var savedData interface{}
var updatedData interface{}
var deletedData []string

type dummyRepository struct {
}

func (s dummyRepository) Create(collectionName string, data interface{}) error {
	savedData = data

	return nil
}

func (s dummyRepository) Read(collectionName string, id string, result *interface{}) error {
	switch id {
	case idFixture:
		*result = map[string]interface{}{
			"ID":   idFixture,
			"Data": dataValueFixture,
			"NestedData": map[string]interface{}{
				"Data":                    nestedDataValueFixture,
				"DeeplyNestedIndexedData": deeplyNestedIDFixture,
			},
			"NestedIndexedData": nestedIDFixture,
			"MappedIndexedData": map[string]string{mapIndexFixture: deeplyNestedIDFixture},
			"SlicedIndexedData": []string{deeplyNestedIDFixture},
			"MappedData":        map[string]string{mapIndexFixture: dataValueFixture},
		}
	case nestedIDFixture:
		*result = map[string]interface{}{
			"ID":                      nestedIDFixture,
			"Data":                    nestedIndexedDataValueFixture,
			"DeeplyNestedIndexedData": deeplyNestedIDFixture,
		}
	case deeplyNestedIDFixture:
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
