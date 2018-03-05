package service

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"regexp"
	"reflect"
)

type EntityDefinitions map[string]reflect.Type

func (e EntityDefinitions) validate() error {
	for _, t := range e {
		if !CanReference(t) {
			return fmt.Errorf("all entities have to be referenceable by having an ID string field but %q does not", t.Name())
		}
	}

	return nil
}

type StorageService struct {
	entityDefinitions EntityDefinitions
	idGenerator       IDGenerator
	repository        Repository
}

func NewStorageService(entityDefinitions EntityDefinitions, idGenerator IDGenerator, repository Repository) (StorageService, error) {
	err := entityDefinitions.validate()
	if err != nil {
		return StorageService{}, err
	}

	return StorageService{
		entityDefinitions: entityDefinitions,
		idGenerator:       idGenerator,
		repository:        repository,
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

	t, err := s.detectType(r)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotFound)

		return
	}

	indexRegex := regexp.MustCompile("^.*/([^/]+)$")
	index := string(indexRegex.ReplaceAll([]byte(r.URL.Path), []byte("$1")))

	reference, err := GetReference(t)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = s.repository.Read(t.Name(), index, &reference)
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

	t, err := s.detectType(r)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotFound)

		return
	}

	reference, err := GetReference(t)
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

	err = AssertExistingReferences(s.repository, reference, t)
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

	err = s.repository.Save(t.Name(), reference)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	rw.Write(response)
}

func (s StorageService) detectType(r *http.Request) (reflect.Type, error) {
	entityNameRegex := regexp.MustCompile("^/(.*)/[^/]+$")

	if !entityNameRegex.Match([]byte(r.URL.Path)) {
		entityNameRegex = regexp.MustCompile("^/(.*)$")
	}

	entityName := string(entityNameRegex.ReplaceAll([]byte(r.URL.Path), []byte("$1")))

	t, ok := s.entityDefinitions[entityName]
	if !ok {
		return nil, fmt.Errorf("entity %s is not defined", entityName)
	}

	return t, nil
}
