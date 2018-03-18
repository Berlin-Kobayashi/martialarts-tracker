package service

import (
	"testing"
	"reflect"
	"github.com/DanShu93/martialarts-tracker/query"
)

func TestAssertExistingResource(t *testing.T) {
	input := createReferenceFixture()

	err := AssertExistingResource(dummyRepository{}, input, reflect.TypeOf(indexedData{}))

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestAssertExistingResourceForMissingReference(t *testing.T) {
	input := createReferenceFixture()
	input["SlicedIndexedData"] = []string{missingIDFixture}
	err := AssertExistingResource(dummyRepository{}, input, reflect.TypeOf(indexedData{}))

	if err == nil {
		t.Error("Expected error")
	}
}

func TestAssertExistingReferencesForMissingResource(t *testing.T) {
	input := createReferenceFixture()
	input["ID"] = missingIDFixture
	err := AssertExistingReferences(dummyRepository{}, input, reflect.TypeOf(indexedData{}))

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
		"SlicedIndexedData": reflect.New(reflect.TypeOf([]string{})).Interface(),
	}

	actual, err := GetReference(input)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result\n\n Actual: %+v\n\n Expected: %+v", actual, expected)
	}
}

func TestGetReferenceUnsupportedFieldType(t *testing.T) {
	input := reflect.TypeOf(unsupportedFieldMap{})

	_, err := GetReference(input)

	if err == nil {
		t.Error("Expected error")
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
		"SlicedIndexedData": []string{deeplyNestedIDFixture},
	}
}

func TestCanReference(t *testing.T) {
	cases := []struct {
		t            reflect.Type
		canReference bool
		name         string
	}{
		{
			reflect.TypeOf(deeplyNestedIndexedData{}),
			true,
			"Legal",
		},
		{
			reflect.TypeOf(1),
			false,
			"WrongType",
		},
		{
			reflect.TypeOf(struct{}{}),
			false,
			"NoID",
		},
		{
			reflect.TypeOf(struct{ ID int }{}),
			false,
			"WrongIDType",
		},
	}
	for _, c := range cases {
		if CanReference(c.t) != c.canReference {
			t.Errorf("Does not return expected result for %q type", c.name)
		}
	}
}

func TestDerefence(t *testing.T) {
	expected := indexedDataFixture
	actual := reflect.New(reflect.TypeOf(indexedData{})).Interface()

	err := Derefence(dummyRepository{}, createReferenceFixture(), &actual)
	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if !reflect.DeepEqual(actual, &expected) {
		t.Errorf("Unexpected result\n\n Actual: %+v\n\n Expected: %+v", actual, expected)
	}
}

func TestGetReferencedBy(t *testing.T) {
	expected := map[string]interface{}{
		"indexedData": []interface{}{createReferenceFixture()},
	}

	actual, err := GetReferencedBy(dummyRepository{}, nestedIDFixture, reflect.TypeOf(nestedIndexedData{}), []reflect.Type{
		reflect.TypeOf(indexedData{}),
		reflect.TypeOf(nestedIndexedData{}),
		reflect.TypeOf(deeplyNestedIndexedData{}),
		reflect.TypeOf(nestedData{}),
	})

	expectedQuery := query.Query{Q: map[string]query.FieldQuery{"NestedIndexedData": {Kind: query.KindContains, Values: []interface{}{nestedIDFixture}}}}

	if err != nil {
		t.Fatalf("Unexpected error %q", err)
	}

	if !reflect.DeepEqual(queriedData, expectedQuery) {
		t.Errorf("Unexpected query.\nActual: %+v\nExpected %+v", queriedData, expectedQuery)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result.\nActual: %+v\nExpected %+v", actual, expected)
	}
}
