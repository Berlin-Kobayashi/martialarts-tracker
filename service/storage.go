package service

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"regexp"
	"github.com/DanShu93/martialarts-tracker/storage"
)

type EntityDefinitions map[string]EntityDefinition

type EntityDefinition struct {
	Entity     interface{}
	Repository Repository
}

type StorageService struct {
	EntityDefinitions EntityDefinitions
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

	indexRegex := regexp.MustCompile("^/[^/]+/([^/]+)$")
	index := string(indexRegex.ReplaceAll([]byte(r.URL.Path), []byte("$1")))

	entityDefinition, err := s.detectEntityDefinition(r)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusNotFound)

		return
	}

	entity := entityDefinition.Entity
	err = entityDefinition.Repository.Read(index, entity)

	if err != nil {
		fmt.Println(err)
		switch err {
		case storage.NotFound:
			rw.WriteHeader(http.StatusNotFound)
		case storage.Invalid:
			rw.WriteHeader(http.StatusInternalServerError)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	response, err := json.Marshal(entity)
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

	entity := entityDefinition.Entity
	err = json.Unmarshal(content, entity)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)

		return
	}

	entity = entityDefinition.Entity
	err = entityDefinition.Repository.Save(entity)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(entity)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	rw.Write(response)
}

func (s StorageService) detectEntityDefinition(r *http.Request) (EntityDefinition, error) {
	entityNameRegex := regexp.MustCompile("^/([^/]+).*")

	entityName := string(entityNameRegex.ReplaceAll([]byte(r.URL.Path), []byte("$1")))

	entityType, ok := s.EntityDefinitions[entityName]
	if !ok {
		return EntityDefinition{}, fmt.Errorf("entity %s is not defined", entityName)
	}

	return entityType, nil
}
