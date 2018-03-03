package service

import (
	"reflect"
	"errors"
)

func GetReferencingEntity(t reflect.Type) (interface{}, error) {
	var referencingEntity interface{}

	switch t.Kind() {
	case reflect.Struct:

	default:
		return nil, errors.New("could not get reference entity for non struct entity")
	}

	return referencingEntity, nil
}
