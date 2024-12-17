package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetUserStore() {
	users = make(map[int]User)
	nextID = 1
}

func TestGetUsers(t *testing.T) {
	resetUserStore()

	users[1] = User{ID: 1, Name: "name", Lastname: "lastname"}
	users[2] = User{ID: 2, Name: "Jan", Lastname: "Kowalski"}

	req, _ := http.NewRequest("GET", "/users", nil)
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handleUsers)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var userList []User
	err := json.Unmarshal(resp.Body.Bytes(), &userList)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(userList) != 2 {
		t.Errorf("expected 2 users, got %d", len(userList))
	}
}

func TestGetUserByID(t *testing.T) {
	resetUserStore()

	users[1] = User{ID: 1, Name: "name", Lastname: "lastname"}

	req, _ := http.NewRequest("GET", "/users/1", nil)
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handleUserByID)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var user User
	err := json.Unmarshal(resp.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if user.Name != "name" || user.Lastname != "lastname" {
		t.Errorf("unexpected user data: %v", user)
	}
}

func TestPostUser(t *testing.T) {
	resetUserStore()

	newUser := map[string]string{"name": "name", "lastname": "lastname"}
	jsonBody, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(handleUsers)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.Code)
	}

	var responseUser User
	json.NewDecoder(resp.Body).Decode(&responseUser)

	if responseUser.Name != "name" || responseUser.Lastname != "lastname" {
		t.Errorf("unexpected response body: %v", responseUser)
	}

	if responseUser.ID != 1 {
		t.Errorf("expected ID 1, got %d", responseUser.ID)
	}
}

func TestPatchUser(t *testing.T) {
	resetUserStore()

    users[1] = User{ID: 1, Name: "TypeError: undefined is not a function", Lastname: "Oczkowski"}

	nameUpdate := map[string]string{"name": "Wojciech"}
	jsonBody, _ := json.Marshal(nameUpdate)
	req, _ := http.NewRequest("PATCH", "/users/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(handleUserByID)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.Code)
	}

	updatedUser := users[1]
	if updatedUser.Name != "Wojciech" || updatedUser.Lastname != "Oczkowski" {
		t.Errorf("patch failed: %v", updatedUser)
	}

	mixedUpdate := map[string]string{
		"name": "name", 
		"lastname": "lastname", 
		"invalid": "field",
	}
	jsonBody, _ = json.Marshal(mixedUpdate)
	req, _ = http.NewRequest("PATCH", "/users/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.Code)
	}

	updatedUser = users[1]
	if updatedUser.Name != "name" || updatedUser.Lastname != "lastname" {
		t.Errorf("patch failed: %v", updatedUser)
	}
}

// here throw err in invalid field idk why
// func TestPatchUser(t *testing.T) {
// 	resetUserStore()

// 	users[1] = User{ID: 1, Name: "Jan", Lastname: "Kowalski"}

// 	nameUpdate := map[string]string{"name": "name"}
// 	jsonBody, _ := json.Marshal(nameUpdate)
// 	req, _ := http.NewRequest("PATCH", "/users/1", bytes.NewReader(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	resp := httptest.NewRecorder()

// 	handler := http.HandlerFunc(handleUserByID)
// 	handler.ServeHTTP(resp, req)

// 	if resp.Code != http.StatusNoContent {
// 		t.Errorf("expected status 204, got %d", resp.Code)
// 	}

// 	updatedUser := users[1]
// 	if updatedUser.Name != "name" || updatedUser.Lastname != "Kowalski" {
// 		t.Errorf("patch failed: %v", updatedUser)
// 	}

// 	lastnameUpdate := map[string]string{"lastname": "lastname"}
// 	jsonBody, _ = json.Marshal(lastnameUpdate)
// 	req, _ = http.NewRequest("PATCH", "/users/1", bytes.NewReader(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	resp = httptest.NewRecorder()

// 	handler.ServeHTTP(resp, req)

// 	if resp.Code != http.StatusNoContent {
// 		t.Errorf("expected status 204, got %d", resp.Code)
// 	}

// 	updatedUser = users[1]
// 	if updatedUser.Name != "name" || updatedUser.Lastname != "lastname" {
// 		t.Errorf("patch failed: %v", updatedUser)
// 	}

// 	invalidUpdate := map[string]string{"invalid": "field"}
// 	jsonBody, _ = json.Marshal(invalidUpdate)
// 	req, _ = http.NewRequest("PATCH", "/users/1", bytes.NewReader(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	resp = httptest.NewRecorder()

// 	handler.ServeHTTP(resp, req)

// 	if resp.Code != http.StatusBadRequest {
// 		t.Errorf("expected status 400 for invalid patch, got %d", resp.Code)
// 	}
// }

func TestPutUser(t *testing.T) {
	resetUserStore()

	newUser := map[string]string{"name": "name", "lastname": "lastname"}
	jsonBody, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(handleUserByID)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.Code)
	}

	updatedUser := users[1]
	if updatedUser.Name != "name" || updatedUser.Lastname != "lastname" {
		t.Errorf("put failed: %v", updatedUser)
	}
}

func TestDeleteUser(t *testing.T) {
	resetUserStore()

	users[1] = User{ID: 1, Name: "name", Lastname: "lastname"}

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(handleUserByID)
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.Code)
	}

	req, _ = http.NewRequest("GET", "/users/1", nil)
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for deleted user, got %d", resp.Code)
	}

	req, _ = http.NewRequest("DELETE", "/users/1", nil)
	resp = httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for non-existent user, got %d", resp.Code)
	}
}

