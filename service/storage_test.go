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

var storageService = StorageService{
	entityDefinitions: EntityDefinitions{
		entityName: {
			T: reflect.TypeOf(deeplyNestedIndexedData{}),
			R: dummyRepository{},
		},
	},
}

func TestStorageService_ServeHTTPGET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", entityName, deeplyNestedIDFixture), nil)
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	expectedBody := getDataFixtureJSON(t)

	content, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != expectedBody {
		t.Errorf("Does not produce expected response. Actual %s Expected %s", content, expectedBody)
	}
}

func TestStorageService_ServeHTTPPOST(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s", entityName), bytes.NewReader([]byte(getDataFixtureJSON(t))))
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	expected := map[string]interface{}{
		"ID":   deeplyNestedIDFixture,
		"Data": deeplyNestedDataValueFixture,
	}

	if !reflect.DeepEqual(recordedData, expected) {
		t.Errorf("Does not save expected data. Actual %v Expected %v", recordedData, deeplyNestedIndexedDataFixture)
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

func getDataFixtureJSON(t *testing.T) string {
	fixtureJSON, err := json.Marshal(deeplyNestedIndexedDataFixture)

	if err != nil {
		t.Errorf("Could not create data fixture.")
	}

	return string(fixtureJSON)
}
