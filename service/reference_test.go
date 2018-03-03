package service

import (
	"testing"
	"reflect"
)

func TestGetReferencingEntity(t *testing.T) {
	input := struct {
		ID   string
		Data string
		Nested struct {
			ID         string
			NestedData string
		}
	}{
		ID:   "myID",
		Data: "myData",
		Nested: struct {
			ID         string
			NestedData string
		}{
			ID:         "myNestedID",
			NestedData: "myNestedData",
		},
	}

	expected := struct {
		ID   string
		Data string
		Nested struct {
			ID string
		}
	}{
		ID:   "myID",
		Data: "myData",
		Nested: struct {
			ID string
		}{
			ID: "myNestedID",
		},
	}

	actual, err := GetReferencingEntity(reflect.TypeOf(input))

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if actual != expected {
		t.Errorf("Unexpected result %s", actual)
	}
}

func TestGetReferencingEntityForNonStruct(t *testing.T) {
	input := 1

	_, err := GetReferencingEntity(reflect.TypeOf(input))

	if err == nil {
		t.Error("Expected error ")
	}
}
