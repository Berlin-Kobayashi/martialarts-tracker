package service

import (
	"reflect"
)

const idFieldName = "ID"

func GetReference(t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.Struct:
		result := make(map[string]interface{}, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			property := t.Field(i)
			if property.Type.Kind() == reflect.Struct {
				if CanReference(property.Type) {
					result[property.Name] = reflect.New(reflect.TypeOf("")).Interface()
				} else {
					referencingEntity, err := GetReference(property.Type)
					if err != nil {
						return nil, err
					}
					result[property.Name] = referencingEntity
				}
			} else {
				referencingEntity, err := GetReference(property.Type)
				if err != nil {
					return nil, err
				}
				result[property.Name] = referencingEntity
			}
		}

		return result, nil
	case reflect.Map:
		property := t.Elem()
		if property.Kind() == reflect.Struct {
			if CanReference(property) {
				return reflect.New(reflect.TypeOf(map[string]string{})).Interface(), nil
			} else {
				return GetReference(property)
			}
		} else {
			referencingEntity, err := GetReference(property)
			if err != nil {
				return nil, err
			}

			return reflect.New(reflect.MapOf(t.Key(), reflect.TypeOf(referencingEntity))).Interface(), nil
		}
	case reflect.Slice:
		property := t.Elem()
		if property.Kind() == reflect.Struct {
			if CanReference(property) {
				return reflect.New(reflect.TypeOf([]string{})).Interface(), nil
			} else {
				return GetReference(property)
			}
		} else {
			return GetReference(property)
		}
	}

	return reflect.New(t).Interface(), nil
}

func Derefence(repository Repository, index string, result *interface{}) error {
	t := reflect.TypeOf(result).Elem()
	reference, err := GetReference(t)
	if err != nil {

		return err
	}

	err = repository.Read(t.Name(), index, &reference)
	if err != nil {
		return err
	}

	//TODO implement

	*result = reference

	return nil
}

func AssertExistingResource(repository Repository, entity interface{}, t reflect.Type) error {
	return assertExistingResourceRecursively(repository, entity, t, true)
}

func AssertExistingReferences(repository Repository, entity interface{}, t reflect.Type) error {
	return assertExistingResourceRecursively(repository, entity, t, false)
}

func assertExistingResourceRecursively(repository Repository, entity interface{}, t reflect.Type, checkRoot bool) error {
	v := reflect.ValueOf(entity)
	switch t.Kind() {
	case reflect.Struct:
		if checkRoot && CanReference(t) {
			id := ""
			if v.Kind() == reflect.String {
				id = v.Interface().(string)
			} else {
				id = v.MapIndex(reflect.ValueOf(idFieldName)).Interface().(string)
			}

			propertyValue := reflect.New(t).Interface()
			if err := repository.Read(t.Name(), id, &propertyValue); err != nil {
				return err
			}
		}

		if v.Kind() != reflect.String {
			for i := 0; i < t.NumField(); i++ {
				fieldValue := v.MapIndex(reflect.ValueOf(t.Field(i).Name))
				err := assertExistingResourceRecursively(repository, fieldValue.Interface(), t.Field(i).Type, true)
				if err != nil {
					return err
				}
			}
		}

		return nil
	case reflect.Map:
		for _, k := range v.MapKeys() {
			err := assertExistingResourceRecursively(repository, v.MapIndex(k).Interface(), t.Elem(), true)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			err := assertExistingResourceRecursively(repository, v.Index(i).Interface(), t.Elem(), true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CanReference(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}

	idField, hasID := t.FieldByName(idFieldName)

	return hasID && idField.Type.Kind() == reflect.String
}
