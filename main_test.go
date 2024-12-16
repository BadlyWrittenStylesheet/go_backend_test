package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsersEmpty(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users", nil)
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handleUsers)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "[]" {
		t.Errorf("expected empty list, got %s", body)
	}
}

func TestPostUser(t *testing.T) {
	newUser := map[string]string{"name": "John", "lastname": "Doe"}
	jsonBody, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(handleUsers)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.Code)
	}

	var responseUser map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseUser)
	if responseUser["name"] != "John" || responseUser["lastname"] != "Doe" {
		t.Errorf("unexpected response body: %v", responseUser)
	}
}

func TestGetUserByID(t *testing.T) {
	newUser := map[string]string{"name": "Alice", "lastname": "Smith"}
	jsonBody, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(handleUsers)
	handler.ServeHTTP(resp, req)

	var createdUser map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&createdUser)
	id := int(createdUser["id"].(float64))

	req, _ = http.NewRequest("GET", "/users/"+string(rune(id)), nil)
	resp = httptest.NewRecorder()

	handler = http.HandlerFunc(handleUserByID)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var fetchedUser map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&fetchedUser)
	if fetchedUser["name"] != "Alice" || fetchedUser["lastname"] != "Smith" {
		t.Errorf("unexpected user data: %v", fetchedUser)
	}
}

func TestDeleteUser(t *testing.T) {
	newUser := map[string]string{"name": "Daisy", "lastname": "Green"}
	jsonBody, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(handleUsers)
	handler.ServeHTTP(resp, req)

	var createdUser map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&createdUser)
	id := int(createdUser["id"].(float64))

	resp = httptest.NewRecorder()

	handler = http.HandlerFunc(handleUserByID)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.Code)
	}

	req, _ = http.NewRequest("GET", "/users/"+string(rune(id)), nil)
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for non-existent user, got %d", resp.Code)
	}
}


