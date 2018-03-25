package service

import (
	"testing"
	"reflect"
	"github.com/DanShu93/martialarts-tracker/query"
)

func TestAssertExistingResource(t *testing.T) {
	input := createReferencingFixture()

	err := AssertExistingResource(dummyRepository{}, input, reflect.TypeOf(referencingData{}))

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}
}

func TestAssertExistingResourceForMissingReference(t *testing.T) {
	input := createReferencingFixture()
	input["SlicedIndexedData"] = []string{missingIDFixture}
	err := AssertExistingResource(dummyRepository{}, input, reflect.TypeOf(referencingData{}))

	if err == nil {
		t.Error("Expected error")
	}
}

func TestAssertExistingReferencesForMissingResource(t *testing.T) {
	input := createReferencingFixture()
	input["ID"] = missingIDFixture
	err := AssertExistingReferences(dummyRepository{}, input, reflect.TypeOf(referencingData{}))

	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}
}

func TestGetReference(t *testing.T) {
	input := reflect.TypeOf(referencingDataFixture)

	expected := map[string]interface{}{
		"ID":   reflect.New(reflect.TypeOf("")).Interface(),
		"Data": reflect.New(reflect.TypeOf("")).Interface(),
		"NestedData" : reflect.New(reflect.TypeOf(nestedData{})).Interface(),
		"References": map[string]interface{}{
			"Single":   reflect.New(reflect.TypeOf("")).Interface(),
			"Multiple": reflect.New(reflect.TypeOf([]string{})).Interface(),
		},
	}

	actual, err := GetReference(input)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result\n\n Actual: %+v\n\n Expected: %+v", actual, expected)
	}
}

func TestCanBeReferenced(t *testing.T) {
	cases := []struct {
		t               reflect.Type
		canBeReferenced bool
		name            string
	}{
		{
			reflect.TypeOf(referencedData{}),
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
		if CanBeReferenced(c.t) != c.canBeReferenced {
			t.Errorf("Does not return expected result for %q", c.name)
		}
	}
}

func TestHasReferences(t *testing.T) {
	cases := []struct {
		t             reflect.Type
		hasReferences bool
		name          string
	}{
		{
			reflect.TypeOf(referencingData{}),
			true,
			"Legal",
		},
		{
			reflect.TypeOf(referencedData{}),
			false,
			"NoReferences",
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
		if HasReferences(c.t) != c.hasReferences {
			t.Errorf("Does not return expected result for %q", c.name)
		}
	}
}

func TestDerefence(t *testing.T) {
	expected := referencingDataFixture
	actual := reflect.New(reflect.TypeOf(referencingData{})).Interface()

	err := Derefence(dummyRepository{}, createReferencingFixture(), &actual)
	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if !reflect.DeepEqual(actual, &expected) {
		t.Errorf("Unexpected result\n\n Actual: %+v\n\n Expected: %+v", actual, expected)
	}
}

func TestGetReferencedBy(t *testing.T) {
	queriedData = []query.Query{}
	expected := map[string]interface{}{
		"referencingData": []interface{}{createReferencingFixture()},
	}

	actual, err := GetReferencedBy(dummyRepository{}, referencedIDFixture, reflect.TypeOf(referencedData{}), []reflect.Type{
		reflect.TypeOf(referencingData{}),
		reflect.TypeOf(referencedData{}),
		reflect.TypeOf(nestedData{}),
	})

	expectedQuery := []query.Query{
		{Q: map[string]query.FieldQuery{"References.Single": {Kind: query.KindAnd, Values: []interface{}{referencedIDFixture}}}},
		{Q: map[string]query.FieldQuery{"References.Multiple": {Kind: query.KindContains, Values: []interface{}{referencedIDFixture}}}},
	}

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

func createReferencingFixture() map[string]interface{} {
	return map[string]interface{}{
		"ID":   referencingIDFixture,
		"Data": referencingValueFixture,
		"NestedData" : nestedDataFixture,
		"References": map[string]interface{}{
			"Single":   referencedIDFixture,
			"Multiple": []string{referencedIDFixture},
		},
	}
}
