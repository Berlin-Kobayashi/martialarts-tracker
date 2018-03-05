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

func (s dummyRepository) Save(data interface{}) error {
	recordedData = data

	return nil
}

func (s dummyRepository) Read(id string, result *interface{}) error {
		switch id {
		case idFixture:
			*result = indexedDataFixture
		case nestedIDFixture:
			*result = nestedIDFixture
		case deeplyNestedIDFixture:
			*result = deeplyNestedIDFixture
		default:
			return storage.NotFound
		}

	return nil
}
