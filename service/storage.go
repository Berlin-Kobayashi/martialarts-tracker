package service

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"regexp"
	"reflect"
)

const ActionExpand = "expand"

var entityNameRegex = regexp.MustCompile("^/([^/]+)/.*$")
var indexRegex = regexp.MustCompile("^.*/([^/]+)$")
var actionRegex = regexp.MustCompile("^/[^/]+/([^/]+)$")
var indexedActionRegex = regexp.MustCompile("^/[^/]+/([^/]+)/[^/]+$")
var indexedEntityNameRegex = regexp.MustCompile("^/([^/]+)$")

var pathRegex = regexp.MustCompilePOSIX("^/[^/]+/?[^/]*/?[^/]*$")

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
	rw.Header().Add("Content-Type", "application/json")

	t, err := s.getType(r)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotFound)

		return
	}

	if !pathRegex.Match([]byte(r.URL.Path)) {
		rw.WriteHeader(http.StatusNotFound)
	}

	index := s.getIndex(r)
	action := s.getAction(r)

	switch r.Method {
	case http.MethodGet:
		switch action {
		case ActionExpand:
			s.expand(rw, r, t, index)
		default:
			s.get(rw, r, t, index)
		}
	case http.MethodPost:
		s.post(rw, r, t)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

	rw.WriteHeader(http.StatusNotFound)
}

func (s StorageService) get(rw http.ResponseWriter, r *http.Request, t reflect.Type, index string) {
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

func (s StorageService) expand(rw http.ResponseWriter, r *http.Request, t reflect.Type, index string) {
	reference, err := GetReference(t)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	result := reflect.New(t).Interface()

	err = Derefence(s.repository, index, reference, &result)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotFound)

		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	rw.Write(response)
}

func (s StorageService) post(rw http.ResponseWriter, r *http.Request, t reflect.Type) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
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

func (s StorageService) getAction(r *http.Request) string {
	regex := actionRegex

	if !regex.Match([]byte(r.URL.Path)) {
		regex = indexedActionRegex
	}
	if !regex.Match([]byte(r.URL.Path)) {
		return ""
	}

	return string(regex.ReplaceAll([]byte(r.URL.Path), []byte("$1")))
}

func (s StorageService) getIndex(r *http.Request) string {
	if !indexRegex.Match([]byte(r.URL.Path)) {
		return ""
	}

	return string(indexRegex.ReplaceAll([]byte(r.URL.Path), []byte("$1")))
}

func (s StorageService) getType(r *http.Request) (reflect.Type, error) {
	regex := entityNameRegex

	if !regex.Match([]byte(r.URL.Path)) {
		regex = indexedEntityNameRegex
	}

	entityName := string(regex.ReplaceAll([]byte(r.URL.Path), []byte("$1")))

	t, ok := s.entityDefinitions[entityName]
	if !ok {
		return nil, fmt.Errorf("entity %s is not defined", entityName)
	}

	return t, nil
}
