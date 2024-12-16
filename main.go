package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID int          `json:"id"`
	Name string     `json:"name"`
	Lastname string `json:"lastname"`
}

var (
	users = make(map[int]User)
	nextID = 1
	usersLock sync.Mutex
)

func main() {
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/", handleUserByID)
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleGetUsers(w, r)
		return
	} else if r.Method == http.MethodPost {
		handleCreateUser(w, r)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func handleUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/users/"):])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetUserByID(w, id)
	case http.MethodPatch:
		handlePatchUser(w, r, id)
	case http.MethodPut:
		handlePutUser(w, r, id)
	case http.MethodDelete:
		handleDeleteUser(w, id)
	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	usersLock.Lock()
	defer usersLock.Unlock()

	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}

func handleGetUserByID(w http.ResponseWriter, id int) {
	usersLock.Lock()
	defer usersLock.Unlock()

	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	usersLock.Lock()
	defer usersLock.Unlock()

	newUser.ID = nextID
	nextID += 1
	users[newUser.ID] = newUser

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func handlePatchUser(w http.ResponseWriter, r *http.Request, id int) {
	var updates map[string]string
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	usersLock.Lock()
	defer usersLock.Unlock()

	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if name, ok := updates["name"]; ok {
		user.Name = name
	}
	if lastname, ok := updates["lastname"]; ok {
		user.Lastname = lastname
	}

	users[id] = user
	w.WriteHeader(http.StatusNoContent)
}

func handlePutUser(w http.ResponseWriter, r *http.Request, id int) {
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	updatedUser.ID = id
	users[id] = updatedUser
	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteUser(w http.ResponseWriter, id int) {
	usersLock.Lock()
	defer usersLock.Unlock()

	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}

