package service

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"regexp"
	"reflect"
)

type EntityDefinitions map[string]EntityDefinition

func (e EntityDefinitions) createEntityStorage() (entityStorage, error) {
	es := make(entityStorage, len(e))
	for _, v := range e {
		t := v.T
		if t.Kind() != reflect.Struct {
			return entityStorage{}, fmt.Errorf("all entities have to be struct but %q is a %q", t.Name(), t.Kind())
		}

		if idField, hasID := t.FieldByName(idFieldName); !hasID || idField.Type.Kind() != reflect.String {
			return entityStorage{}, fmt.Errorf("all entities have to have an ID string field but %q does not", t.Name())
		}

		es[t] = v.R
	}

	return es, nil
}

type EntityDefinition struct {
	R Repository
	T reflect.Type
}

type StorageService struct {
	entityDefinitions EntityDefinitions
	entityStorage     entityStorage
	idGenerator       IDGenerator
}

func NewStorageService(entityDefinitions EntityDefinitions, idGenerator IDGenerator) (StorageService, error) {
	entityStorage, err := entityDefinitions.createEntityStorage()
	if err != nil {
		return StorageService{}, err
	}

	return StorageService{
		entityDefinitions: entityDefinitions,
		entityStorage:     entityStorage,
		idGenerator:       idGenerator,
	}, nil
}

func (s StorageService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.get(rw, r)
	case http.MethodPost:
		s.post(rw, r)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s StorageService) get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	entityDefinition, err := s.detectEntityDefinition(r)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotFound)

		return
	}

	indexRegex := regexp.MustCompile("^.*/([^/]+)$")
	index := string(indexRegex.ReplaceAll([]byte(r.URL.Path), []byte("$1")))

	reference, err := GetReference(entityDefinition.T)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = entityDefinition.R.Read(index, &reference)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotFound)

		return
	}

	response, err := json.Marshal(reference)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	rw.Write(response)
}

func (s StorageService) post(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	entityDefinition, err := s.detectEntityDefinition(r)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotFound)

		return
	}

	reference, err := GetReference(entityDefinition.T)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.Unmarshal(content, &reference)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)

		return
	}

	id := s.idGenerator.Generate()
	reference.(map[string]interface{})[idFieldName] = id

	err = s.entityStorage.AssertExistingReferences(reference, entityDefinition.T)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)

		return
	}

	response, err := json.Marshal(reference)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = entityDefinition.R.Save(reference)
	if err != nil {
		fmt.Println(err)
		switch err {
		case UnsupportedMethod:
			rw.WriteHeader(http.StatusMethodNotAllowed)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	rw.Write(response)
}

func (s StorageService) detectEntityDefinition(r *http.Request) (EntityDefinition, error) {
	entityNameRegex := regexp.MustCompile("^/(.*)/[^/]+$")

	if !entityNameRegex.Match([]byte(r.URL.Path)) {
		entityNameRegex = regexp.MustCompile("^/(.*)$")
	}

	entityName := string(entityNameRegex.ReplaceAll([]byte(r.URL.Path), []byte("$1")))

	entityType, ok := s.entityDefinitions[entityName]
	if !ok {
		return EntityDefinition{}, fmt.Errorf("entity %s is not defined", entityName)
	}

	return entityType, nil
}
