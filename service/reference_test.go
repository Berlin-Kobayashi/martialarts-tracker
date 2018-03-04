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

func TestEntityStorage_AssertExistingResource(t *testing.T) {
	input := indexedDataFixture
	err := indexedDataStorage.AssertExistingResource(input)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestEntityStorage_AssertExistingResourceForMissingReference(t *testing.T) {
	input := indexedDataFixture

	mappedData := deeplyNestedIndexedDataFixture
	mappedData.ID = "123"
	input.MappedIndexedData = map[string]deeplyNestedIndexedData{mapIndexFixture: mappedData}
	err := indexedDataStorage.AssertExistingResource(input)

	if err == nil {
		t.Error("Expected error")
	}
}

func TestEntityStorage_AssertExistingReferencesForMissingResource(t *testing.T) {
	input := indexedDataFixture
	input.ID = "123"
	err := indexedDataStorage.AssertExistingReferences(input)

	if err != nil {
		t.Errorf("Unexpected error %q", err)
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
		"MappedIndexedData": reflect.New(reflect.TypeOf(map[string]string{})).Interface(),
		"SlicedIndexedData": reflect.New(reflect.TypeOf([]string{})).Interface(),
		"MappedData":        reflect.New(reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(reflect.New(reflect.TypeOf("")).Interface()))).Interface(),
	}

	actual, err := GetReference(input)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result\n\n Actual: %+v\n\n Expected: %+v", actual, expected)
	}
}
