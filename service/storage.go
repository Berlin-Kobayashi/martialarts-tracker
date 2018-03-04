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

func (e EntityDefinitions) createEntityStorage() entityStorage {
	entityStorage := make(entityStorage, len(e))
	for _, v := range e {
		entityStorage[v.T] = v.R
	}

	return entityStorage
}

type EntityDefinition struct {
	R Repository
	T reflect.Type
}

type StorageService struct {
	entityDefinitions EntityDefinitions
	entityStorage     entityStorage
}

func NewStorageService(entityDefinitions EntityDefinitions) StorageService {
	return StorageService{
		entityDefinitions: entityDefinitions,
		entityStorage:     entityDefinitions.createEntityStorage(),
	}
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

	err = s.entityStorage.AssertValidReference(reference)
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

	response, err := json.Marshal(reference)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

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
