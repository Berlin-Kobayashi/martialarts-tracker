package service

import (
	"github.com/DanShu93/martialarts-tracker/storage"
)

type dummyUUIDGenerator struct {
}

func (g dummyUUIDGenerator) Generate() string {
	return uuidV4Fixture
}

var recordedData interface{}

type dummyRepository struct {
}

func (s dummyRepository) Save(collectionName string, data interface{}) error {
	recordedData = data

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
			"ID":                       nestedIDFixture,
			"Data":                     nestedIndexedDataValueFixture,
			"DeeplyNestedIndexedData" : deeplyNestedIDFixture,
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
