package service

import (
	"reflect"
	"fmt"
	"github.com/DanShu93/martialarts-tracker/query"
	"github.com/DanShu93/martialarts-tracker/storage"
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
		return nil, fmt.Errorf("unsupported field type %q", t.Kind())
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

func Derefence(repository Repository, reference, result interface{}) error {
	t := reflect.Indirect(reflect.ValueOf(result)).Type()
	v := reflect.ValueOf(reference)
	res := reflect.Indirect(reflect.ValueOf(result))

	if res.Kind() == reflect.Interface {
		res = res.Elem()
		t = reflect.Indirect(reflect.ValueOf(result)).Elem().Type()

		if res.Kind() == reflect.Ptr {
			res = res.Elem()
			t = reflect.Indirect(reflect.ValueOf(result)).Elem().Elem().Type()
		}
	}

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			fieldValue := v.MapIndex(reflect.ValueOf(t.Field(i).Name)).Elem()
			if CanReference(t.Field(i).Type) {
				subReference, err := GetReference(t.Field(i).Type)
				if err != nil {
					return err
				}

				err = repository.Read(t.Field(i).Type.Name(), fieldValue.Interface().(string), &subReference)
				if err != nil {
					return err
				}

				subResult := reflect.New(t.Field(i).Type).Interface()
				err = Derefence(repository, subReference, &subResult)
				if err != nil {
					return err
				}

				subResultValue := reflect.ValueOf(subResult)
				if subResultValue.Kind() == reflect.Ptr {
					subResultValue = subResultValue.Elem()
				}
				res.Field(i).Set(subResultValue)
			} else {
				switch t.Field(i).Type.Kind() {
				case reflect.Struct, reflect.Slice:
					subResult := reflect.New(t.Field(i).Type).Interface()
					if t.Field(i).Type.Kind() == reflect.Slice {
						subResult = reflect.MakeSlice(t.Field(i).Type, fieldValue.Len(), fieldValue.Cap()).Interface()
					}
					err := Derefence(repository, fieldValue.Interface(), &subResult)
					if err != nil {
						return err
					}

					subResultValue := reflect.ValueOf(subResult)
					if subResultValue.Kind() == reflect.Ptr {
						subResultValue = subResultValue.Elem()
					}
					res.Field(i).Set(subResultValue)
				default:
					res.Field(i).Set(fieldValue)
				}
			}
		}
		return nil
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			fieldValue := v.Index(i)
			if CanReference(t.Elem()) {
				subReference, err := GetReference(t.Elem())
				if err != nil {
					return err
				}
				err = repository.Read(t.Elem().Name(), fieldValue.Interface().(string), &subReference)
				if err != nil {
					return err
				}

				subResult := reflect.New(t.Elem()).Interface()
				err = Derefence(repository, subReference, &subResult)
				if err != nil {
					return err
				}

				subResultValue := reflect.ValueOf(subResult)
				if subResultValue.Kind() == reflect.Ptr {
					subResultValue = subResultValue.Elem()
				}
				res.Index(i).Set(subResultValue)
			} else {
				res.Index(i).Set(fieldValue)
			}
		}
	default:
		*result.(*interface{}) = reference
	}

	return nil
}

func AssertExistingResource(repository Repository, reference interface{}, t reflect.Type) error {
	return assertExistingResourceRecursively(repository, reference, t, true)
}

func AssertExistingReferences(repository Repository, reference interface{}, t reflect.Type) error {
	return assertExistingResourceRecursively(repository, reference, t, false)
}

func assertExistingResourceRecursively(repository Repository, reference interface{}, t reflect.Type, checkRoot bool) error {
	v := reflect.ValueOf(reference)
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

func GetReferencedBy(repository Repository, id string, resourceType reflect.Type, types []reflect.Type) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	for _, t := range types {
		if t.Kind() == reflect.Struct {
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)

				if f.Type == resourceType {
					var references []interface{}
					q := query.Query{Q: map[string]query.FieldQuery{f.Name: {Kind: query.KindContains, Values: []interface{}{id}}}}
					err := repository.ReadAll(t.Name(), q, &references)
					if err != nil {
						if err == storage.NotFound {
							references = []interface{}{}
						} else {
							return nil, err
						}
					}
					result[t.Name()] = references
				}
			}
		}
	}

	return result, nil
}
