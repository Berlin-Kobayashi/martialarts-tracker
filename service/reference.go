package service

import (
	"reflect"
	"fmt"
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
	fmt.Println("A1", res, res.Type(), t)

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			fieldValue := v.MapIndex(reflect.ValueOf(t.Field(i).Name)).Elem()
			if CanReference(t.Field(i).Type) {
				//TODO implement
			} else {
				fmt.Println("B0", t.Field(i).Type.Kind())
				switch t.Field(i).Type.Kind() {
				case reflect.Struct, reflect.Map, reflect.Slice:
					subResult := reflect.New(t.Field(i).Type).Interface()
					switch t.Field(i).Type.Kind() {
					case reflect.Map:
						subResult = reflect.MakeMap(t.Field(i).Type).Interface()
					case reflect.Slice:
						subResult = reflect.MakeSlice(t.Field(i).Type, 0, 0).Interface()
					}

					fmt.Println("B1", reflect.ValueOf(subResult).Type())
					err := Derefence(repository, fieldValue.Interface(), &subResult)
					if err != nil {
						return err
					}

					fmt.Println("B", reflect.ValueOf(subResult).Type())
					fmt.Println("C", subResult)
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
	case reflect.Map:
		for _, k := range v.MapKeys() {
			fieldValue := v.MapIndex(k)
			if CanReference(t.Elem()) {
				//TODO implement
			} else {
				res.SetMapIndex(k, fieldValue)
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			fieldValue := v.Index(i)
			if CanReference(t.Elem()) {
				//TODO implement
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
