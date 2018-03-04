package service

import (
	"testing"
	"reflect"
)

var indexedDataStorage = entityStorage{
	reflect.TypeOf(indexedData{}):             dummyRepository{},
	reflect.TypeOf(nestedIndexedData{}):       dummyRepository{},
	reflect.TypeOf(deeplyNestedIndexedData{}): dummyRepository{},
}

func TestEntityStorage_GetValidReference(t *testing.T) {
	input := indexedDataFixture

	expected := map[string]interface{}{
		"ID":   idFixture,
		"Data": dataValueFixture,
		"NestedData": map[string]interface{}{
			"Data":                    nestedDataValueFixture,
			"DeeplyNestedIndexedData": deeplyNestedIDFixture,
		},
		"NestedIndexedData": nestedIDFixture,
	}

	actual, err := indexedDataStorage.GetValidReference(input)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result %+v", actual)
	}
}

func TestEntityStorage_GetValidReferenceForUnsupportedType(t *testing.T) {
	input := 1

	_, err := indexedDataStorage.GetValidReference(input)

	if err == nil {
		t.Error("Expected error ")
	}
}

func TestGetReference(t *testing.T) {
	input := reflect.TypeOf(indexedDataFixture)

	expected := map[string]interface{}{
		"ID":   reflect.New(reflect.TypeOf("")).Interface(),
		"Data": reflect.New(reflect.TypeOf("")).Interface(),
		"NestedData": map[string]interface{}{
			"Data":                    reflect.New(reflect.TypeOf("")).Interface(),
			"DeeplyNestedIndexedData": reflect.New(reflect.TypeOf("")).Interface(),
		},
		"NestedIndexedData": reflect.New(reflect.TypeOf("")).Interface(),
	}

	actual, err := GetReference(input)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result %+v", actual)
	}
}

func TestGetReferenceForUnsupportedType(t *testing.T) {
	input := 1

	_, err := GetReference(reflect.TypeOf(input))

	if err == nil {
		t.Error("Expected error ")
	}
}
