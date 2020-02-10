package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/redhat-developer/tekton-hub/backend/api/pkg/api"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/models"
)

func TestDatabaseConnection(t *testing.T) {
	if err := models.StartConnection(); err != nil {
		t.Log(err)
		t.Log("Connection error : ", err)
	}
}

func TestGetAllTasksAPI(t *testing.T) {
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetAllTasks)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
}

func TestGetTasksAPI(t *testing.T) {
	req, err := http.NewRequest("GET", "/task/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetTaskWithID)
	q := req.URL.Query()
	q.Add("name", "argocd")
	req.URL.RawQuery = q.Encode()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
}
