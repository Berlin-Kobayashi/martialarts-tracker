package service
//
//import (
//	"testing"
//	"net/http/httptest"
//	"io/ioutil"
//	"github.com/DanShu93/martialarts-tracker/entity"
//	"fmt"
//	"encoding/json"
//	"bytes"
//	"reflect"
//	"net/http"
//)
//
//var entityName = "training-unit"
//
//var storageService = StorageService{
//	entityDefinitions: EntityDefinitions{
//		entityName: {
//			T: reflect.TypeOf(entity.TrainingUnit{}),
//			R: dummyRepository{},
//		},
//	},
//}
//
//func TestStorageService_ServeHTTPGET(t *testing.T) {
//	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", entityName, trainingUnitFixture.ID), nil)
//	w := httptest.NewRecorder()
//	storageService.ServeHTTP(w, req)
//
//	expectedBody := getTrainingUnitFixtureJSON(t)
//
//	content, err := ioutil.ReadAll(w.Result().Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if string(content) != expectedBody {
//		t.Errorf("Does not produce expected response. Actual %q Expected %q", content, expectedBody)
//	}
//}
//
//func TestStorageService_ServeHTTPPOST(t *testing.T) {
//	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s", entityName), bytes.NewReader([]byte(getTrainingUnitFixtureJSON(t))))
//	w := httptest.NewRecorder()
//	storageService.ServeHTTP(w, req)
//
//	if !reflect.DeepEqual(recordedData, &trainingUnitFixture) {
//		t.Errorf("Does not save expected data. Actual %v Expected %v", recordedData, trainingUnitFixture)
//	}
//
//	expectedBody := getTrainingUnitFixtureJSON(t)
//
//	content, err := ioutil.ReadAll(w.Result().Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if string(content) != expectedBody {
//		t.Errorf("Does not produce expected response. Actual %q Expected %q", content, expectedBody)
//	}
//}
//
//func TestStorageService_ServeHTTPUnknownMethod(t *testing.T) {
//	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/%s", entityName), nil)
//	w := httptest.NewRecorder()
//	storageService.ServeHTTP(w, req)
//
//	if w.Result().StatusCode != http.StatusMethodNotAllowed {
//		t.Errorf("Does not return proper status for unknown method. expected data Actual %v Expected %v", w.Result().StatusCode, http.StatusMethodNotAllowed)
//	}
//}
//
//func TestStorageService_ServeHTTPGETUnknownEntity(t *testing.T) {
//	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/unknown/%s", trainingUnitFixture.ID), nil)
//	w := httptest.NewRecorder()
//	storageService.ServeHTTP(w, req)
//
//	if w.Result().StatusCode != http.StatusNotFound {
//		t.Errorf("Does not return proper status for unkown entity. expected data Actual %v Expected %v", w.Result().StatusCode, http.StatusMethodNotAllowed)
//	}
//}
//
//func TestStorageService_ServeHTTPPOSTTUnknownEntity(t *testing.T) {
//	req := httptest.NewRequest(http.MethodPost, "/unknown", bytes.NewReader([]byte(getTrainingUnitFixtureJSON(t))))
//	w := httptest.NewRecorder()
//	storageService.ServeHTTP(w, req)
//
//	if w.Result().StatusCode != http.StatusNotFound {
//		t.Errorf("Does not return proper status for unkown entity. expected data Actual %v Expected %v", w.Result().StatusCode, http.StatusMethodNotAllowed)
//	}
//}
//
//func getTrainingUnitFixtureJSON(t *testing.T) string {
//	fixtureJSON, err := json.Marshal(trainingUnitFixture)
//
//	if err != nil {
//		t.Errorf("Could not create training unit fixture.")
//	}
//
//	return string(fixtureJSON)
//}
