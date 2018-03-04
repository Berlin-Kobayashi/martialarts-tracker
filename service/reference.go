package service

import (
	"reflect"
	"errors"
)

const idFieldName = "ID"

type entityStorage map[reflect.Type]Repository

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
				result[property.Name] = reflect.New(property.Type).Interface()
			}
		}

		return result, nil
	}

	return nil, errors.New("could not get reference for non struct entity")
}

// TODO add custom struct tags to overwrite property name
func (e entityStorage) GetValidReference(entity interface{}) (interface{}, error) {
	v := reflect.ValueOf(entity)

	switch v.Type().Kind() {
	case reflect.Struct:
		result := make(map[string]interface{}, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			property := v.Field(i)
			propertyName := v.Type().Field(i).Name
			if property.Kind() == reflect.Struct {
				if idField, hasID := property.Type().FieldByName(idFieldName); hasID && idField.Type.Kind() == reflect.String {
					id := property.FieldByName(idFieldName).String()
					propertyValue := reflect.New(property.Type()).Interface()
					if err := e[property.Type()].Read(id, &propertyValue); err != nil {
						return nil, err
					}

					result[propertyName] = id
				} else {
					referencingEntity, err := e.GetValidReference(property.Interface())
					if err != nil {
						return nil, err
					}
					result[propertyName] = referencingEntity
				}
			} else {
				result[propertyName] = property.Interface()
			}
		}

		return result, nil
	}

	return nil, errors.New("could not get reference for non struct entity")
}
