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
	input := createReferenceFixture()

	err := indexedDataStorage.AssertExistingResource(input, reflect.TypeOf(indexedData{}))

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestEntityStorage_AssertExistingResourceForMissingReference(t *testing.T) {
	input := createReferenceFixture()
	input["MappedIndexedData"] = map[string]string{mapIndexFixture: "123"}
	err := indexedDataStorage.AssertExistingResource(input, reflect.TypeOf(indexedData{}))

	if err == nil {
		t.Error("Expected error")
	}
}

func TestEntityStorage_AssertExistingReferencesForMissingResource(t *testing.T) {
	input := createReferenceFixture()
	input["ID"] = "123"
	err := indexedDataStorage.AssertExistingReferences(input, reflect.TypeOf(indexedData{}))

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

func createReferenceFixture() map[string]interface{} {
	return map[string]interface{}{
		"ID":   idFixture,
		"Data": dataValueFixture,
		"NestedData": map[string]interface{}{
			"Data":                    nestedDataValueFixture,
			"DeeplyNestedIndexedData": deeplyNestedIDFixture,
		},
		"NestedIndexedData": nestedIDFixture,
		"MappedIndexedData": map[string]string{mapIndexFixture:deeplyNestedIDFixture},
		"SlicedIndexedData": []string{deeplyNestedIDFixture},
		"MappedData":        map[string]string{mapIndexFixture: dataValueFixture},
	}
}
