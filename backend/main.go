package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Todo represents a single todo item
type Todo struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

var (
	todos  = make(map[int]Todo)
	nextID = 1
	mu     sync.Mutex
)

func main() {
	// Create API routes group
	http.HandleFunc("/api/todos", todosHandler)
	http.HandleFunc("/api/todos/", todoHandler) // Note the trailing slash for specific ID

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("API Call: %s %s\n", r.Method, r.URL.Path)
	log.Printf("Current todos: %+v\n", todos)

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	switch r.Method {
	case http.MethodGet:
		list := make([]Todo, 0, len(todos))
		for _, todo := range todos {
			list = append(list, todo)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
		log.Printf("GET response: %+v\n", list)

	case http.MethodPost:
		var todo Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		todo.ID = nextID
		nextID++
		todos[todo.ID] = todo
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
		log.Printf("Created todo: %+v\n", todo)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("API Call: %s %s\n", r.Method, r.URL.Path)
	log.Printf("Current todos: %+v\n", todos)

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	idStr := r.URL.Path[len("/api/todos/"):] // Updated path
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		todo, ok := todos[id]
		if !ok {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
		log.Printf("GET response for ID %d: %+v\n", id, todo)

	case http.MethodPut:
		var updatedTodo Todo
		if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, ok := todos[id]; !ok {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		updatedTodo.ID = id // Ensure ID remains the same
		todos[id] = updatedTodo
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedTodo)
		log.Printf("Updated todo %d: %+v\n", id, updatedTodo)

	case http.MethodDelete:
		if _, ok := todos[id]; !ok {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		delete(todos, id)
		w.WriteHeader(http.StatusNoContent)
		log.Printf("Deleted todo %d\n", id)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
