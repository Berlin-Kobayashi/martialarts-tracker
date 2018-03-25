package service

import (
	"reflect"
	"github.com/DanShu93/martialarts-tracker/query"
	"github.com/DanShu93/martialarts-tracker/storage"
	"fmt"
)

const idFieldName = "ID"
const referencesFieldName = "References"

func GetReference(t reflect.Type) (interface{}, error) {
	if !HasReferences(t) {
		return reflect.New(t).Interface(), nil
	}

	result := make(map[string]interface{}, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		property := t.Field(i)
		if property.Name == referencesFieldName {
			referencesMap := make(map[string]interface{}, property.Type.NumField())
			for j := 0; j < property.Type.NumField(); j++ {
				reference := property.Type.Field(j)

				switch reference.Type.Kind() {
				case reflect.Struct:
					if CanBeReferenced(reference.Type) {
						referencesMap[reference.Name] = reflect.New(reflect.TypeOf("")).Interface()
					} else {
						return nil, fmt.Errorf("cannot reference stuct %q", reference.Type.Name())
					}
				case reflect.Slice:
					if CanBeReferenced(reference.Type.Elem()) {
						referencesMap[reference.Name] = reflect.New(reflect.TypeOf([]string{})).Interface()
					} else {
						return nil, fmt.Errorf("cannot reference slice of %q", reference.Type.Name())
					}
				default:
					return nil, fmt.Errorf("cannot reference type %q", reference.Type.Name())
				}

				result[property.Name] = referencesMap
			}
		} else {
			result[property.Name] = reflect.New(property.Type).Interface()
		}
	}

	return result, nil
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

	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.MapIndex(reflect.ValueOf(t.Field(i).Name)).Elem()

		if t.Field(i).Name == referencesFieldName {
			for j := 0; j < t.Field(i).Type.NumField(); j++ {
				referenceValue := fieldValue.MapIndex(reflect.ValueOf(t.Field(i).Type.Field(j).Name)).Elem()
				fmt.Println(referenceValue)

				switch referenceValue.Kind() {
				case reflect.String:
					subReference, err := GetReference(t.Field(i).Type.Field(j).Type)
					if err != nil {
						return err
					}
					err = repository.Read(t.Field(i).Type.Field(j).Type.Name(), referenceValue.Interface().(string), &subReference)
					if err != nil {
						return err
					}

					res.Field(i).Field(j).Set(reflect.ValueOf(subReference))
				case reflect.Slice:
				default:
					return fmt.Errorf("cannot dereference ID %q", fieldValue)
				}

			}
		} else {
			res.Field(i).Set(fieldValue)
		}
	}

	return nil
}

func toStruct(source map[string]interface{}, target interface{}) {
	// TODO implement
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
		if checkRoot && CanBeReferenced(t) {
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

func CanBeReferenced(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}

	idField, hasID := t.FieldByName(idFieldName)

	return hasID && idField.Type.Kind() == reflect.String
}

func HasReferences(t reflect.Type) bool {
	if !CanBeReferenced(t) {
		return false
	}

	referencesField, hasReferences := t.FieldByName(referencesFieldName)

	return hasReferences && referencesField.Type.Kind() == reflect.Struct
}

func GetReferencedBy(repository Repository, id string, resourceType reflect.Type, types []reflect.Type) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	for _, t := range types {
		if t.Kind() == reflect.Struct {
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)

				fType := f.Type
				if fType.Kind() == reflect.Slice {
					fType = fType.Elem()
				}

				if fType == resourceType {
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
