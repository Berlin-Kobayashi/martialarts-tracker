package service

import (
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"bytes"
	"reflect"
	"net/http"
)

var entityName = "data"

var storageService = createStorageService()

func TestNewStorageService(t *testing.T) {
	expected := StorageService{
		entityDefinitions: EntityDefinitions{
			entityName: reflect.TypeOf(deeplyNestedIndexedData{}),
		},
		repository:  dummyRepository{},
		idGenerator: dummyUUIDGenerator{},
	}

	if !reflect.DeepEqual(storageService, expected) {
		t.Errorf("Does not return expected result. Actual %s Expected %s", storageService, expected)
	}
}

func TestNewStorageServiceWrongEntityType(t *testing.T) {
	_, err := NewStorageService(
		EntityDefinitions{
			entityName: reflect.TypeOf(1),
		},
		dummyUUIDGenerator{},
		dummyRepository{},
	)

	if err == nil {
		t.Error("Expected error")
	}
}

func TestNewStorageServiceNoIDEntity(t *testing.T) {
	_, err := NewStorageService(
		EntityDefinitions{
			entityName: reflect.TypeOf(struct{}{}),
		},
		dummyUUIDGenerator{},
		dummyRepository{},
	)

	if err == nil {
		t.Error("Expected error")
	}
}

func TestNewStorageServiceWrongIDTypeEntity(t *testing.T) {
	_, err := NewStorageService(
		EntityDefinitions{
			entityName: reflect.TypeOf(struct{ ID int }{}),
		},
		dummyUUIDGenerator{},
		dummyRepository{},
	)

	if err == nil {
		t.Error("Expected error")
	}
}

func TestStorageService_ServeHTTPGET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", entityName, deeplyNestedIDFixture), nil)
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	expected := map[string]interface{}{
		"ID":   deeplyNestedIDFixture,
		"Data": deeplyNestedDataValueFixture,
	}

	content, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	response := map[string]interface{}{}
	err = json.Unmarshal(content, &response)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Does not produce expected response. Actual %q Expected %q", response, expected)
	}
}

func TestStorageService_ServeHTTPGETExpand(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s/%s", entityName, ActionExpand, deeplyNestedIDFixture), nil)
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	expected := map[string]interface{}{
		"ID":   deeplyNestedIDFixture,
		"Data": deeplyNestedDataValueFixture,
	}

	content, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	response := map[string]interface{}{}
	err = json.Unmarshal(content, &response)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Does not produce expected response. Actual %q Expected %q", response, expected)
	}
}

func TestStorageService_ServeHTTPPOST(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s", entityName), bytes.NewReader([]byte(getDataFixtureJSON(t))))
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	expected := map[string]interface{}{
		"ID":   uuidV4Fixture,
		"Data": deeplyNestedDataValueFixture,
	}

	if !reflect.DeepEqual(savedData, expected) {
		t.Errorf("Does not save expected data. Actual %v Expected %v", savedData, deeplyNestedIndexedDataFixture)
	}

	content, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	response := map[string]interface{}{}
	err = json.Unmarshal(content, &response)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Does not produce expected response. Actual %q Expected %q", response, expected)
	}
}

func TestStorageService_ServeHTTPUnknownMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/%s", entityName), nil)
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Does not return proper status for unknown method. expected data Actual %v Expected %v", w.Result().StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestStorageService_ServeHTTPGETUnknownEntity(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/unknown/%s", deeplyNestedIDFixture), nil)
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Errorf("Does not return proper status for unkown entity. expected data Actual %v Expected %v", w.Result().StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestStorageService_ServeHTTPPOSTTUnknownEntity(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/unknown", bytes.NewReader([]byte(getDataFixtureJSON(t))))
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Errorf("Does not return proper status for unkown entity. expected data Actual %v Expected %v", w.Result().StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestStorageService_ServeHTTPPUT(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/%s/%s", entityName, deeplyNestedIDFixture), bytes.NewReader([]byte(getDataFixtureJSON(t))))
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	expected := map[string]interface{}{
		"ID":   deeplyNestedIDFixture,
		"Data": deeplyNestedDataValueFixture,
	}

	if !reflect.DeepEqual(updatedData, expected) {
		t.Errorf("Does not update expected data. Actual %v Expected %v", updatedData, deeplyNestedIndexedDataFixture)
	}

	content, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	response := map[string]interface{}{}
	err = json.Unmarshal(content, &response)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Does not produce expected response. Actual %q Expected %q", response, expected)
	}
}

func TestStorageService_ServeHTTPPUTUnknownResource(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/%s/%s", entityName, "123"), bytes.NewReader([]byte(getDataFixtureJSON(t))))
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Does not return proper status. expected data Actual %v Expected %v", w.Result().StatusCode, http.StatusBadRequest)
	}
}

func TestStorageService_ServeHTTPPUTNotMatchingIDs(t *testing.T) {
	fixture := deeplyNestedIndexedDataFixture
	fixture.ID = "123"
	fixtureJSON, err := json.Marshal(fixture)
	if err != nil {
		t.Errorf("Could not create data fixture.")
	}

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/%s/%s", entityName, deeplyNestedIDFixture), bytes.NewReader(fixtureJSON))
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Does not return proper status. expected data Actual %v Expected %v", w.Result().StatusCode, http.StatusBadRequest)
	}
}

func TestStorageService_ServeHTTPDELETE(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s/%s", entityName, deeplyNestedIDFixture), bytes.NewReader([]byte(getDataFixtureJSON(t))))
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	expected := []string{deeplyNestedIDFixture}

	if !reflect.DeepEqual(deletedData, expected) {
		t.Errorf("Does not delete expected data. Actual %v Expected %v", deletedData, expected)
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Does not return proper status. expected data Actual %v Expected %v", w.Result().StatusCode, http.StatusOK)
	}

	deletedData = []string{}
}

func TestStorageService_ServeHTTPDELETEUnknownResource(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s/%s", entityName, "123"), bytes.NewReader([]byte(getDataFixtureJSON(t))))
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Errorf("Does not return proper status. expected data Actual %v Expected %v", w.Result().StatusCode, http.StatusNotFound)
	}

	deletedData = []string{}
}

func createStorageService() StorageService {
	storageService, err := NewStorageService(
		EntityDefinitions{
			entityName: reflect.TypeOf(deeplyNestedIndexedData{}),
		},
		dummyUUIDGenerator{},
		dummyRepository{},
	)

	if err != nil {
		panic(err)
	}

	return storageService
}

func getDataFixtureJSON(t *testing.T) string {
	fixtureJSON, err := json.Marshal(deeplyNestedIndexedDataFixture)

	if err != nil {
		t.Errorf("Could not create data fixture.")
	}

	return string(fixtureJSON)
}
