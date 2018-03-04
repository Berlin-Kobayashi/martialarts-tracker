package service

import (
	"reflect"
)

const idFieldName = "ID"

type entityStorage map[reflect.Type]Repository

// TODO add custom struct tags to overwrite property name
func GetReference(t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.Struct:
		result := make(map[string]interface{}, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			property := t.Field(i)
			if property.Type.Kind() == reflect.Struct {
				if idField, hasID := property.Type.FieldByName(idFieldName); hasID && idField.Type.Kind() == reflect.String {
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
			if idField, hasID := property.FieldByName(idFieldName); hasID && idField.Type.Kind() == reflect.String {
				return reflect.New(reflect.TypeOf(map[string]string{})).Interface(), nil
			} else {
				return GetReference(property)
			}
		} else {
			return GetReference(property)
		}
	case reflect.Slice:
		property := t.Elem()
		if property.Kind() == reflect.Struct {
			if idField, hasID := property.FieldByName(idFieldName); hasID && idField.Type.Kind() == reflect.String {
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

func (e entityStorage) AssertExistingResource(entity interface{}) error {
	return e.assertExistingResourceRecursively(entity, true)
}

func (e entityStorage) AssertExistingReferences(entity interface{}) error {
	return e.assertExistingResourceRecursively(entity, false)
}

func (e entityStorage) assertExistingResourceRecursively(entity interface{}, checkRoot bool) error {
	v := reflect.ValueOf(entity)

	switch v.Type().Kind() {
	case reflect.Struct:
		if idField, hasID := v.Type().FieldByName(idFieldName); checkRoot && hasID && idField.Type.Kind() == reflect.String {
			id := v.FieldByName(idFieldName).String()
			propertyValue := reflect.New(v.Type()).Interface()
			if err := e[v.Type()].Read(id, &propertyValue); err != nil {
				return err
			}
		}
		for i := 0; i < v.NumField(); i++ {
			err := e.assertExistingResourceRecursively(v.Field(i).Interface(), true)
			if err != nil {
				return err
			}
		}

		return nil

	case reflect.Map:
		for _, k := range v.MapKeys() {
			err := e.assertExistingResourceRecursively(v.MapIndex(k).Interface(), true)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			err := e.assertExistingResourceRecursively(v.Index(i).Interface(), true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
