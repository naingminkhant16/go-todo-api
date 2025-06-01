package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Task struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var tasks = []Task{
	{Id: 1, Name: "To learn Go"},
	{Id: 2, Name: "To read book"},
	{Id: 3, Name: "To watch movie"},
	{Id: 4, Name: "To go to gym"},
}

func main() {
	serveHttp(registerRoutes())
}

func serveHttp(router *mux.Router) {
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}

func registerRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", getAllTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", getTaskById).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", deleteTaskById).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateById).Methods("PUT")
	return router
}

func generateNewId() int {
	return tasks[len(tasks)-1].Id + 1
}

func updateById(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	idStr := mux.Vars(request)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var task Task
	err = json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var found bool = false
	for index, t := range tasks {
		if t.Id != id {
			continue
		}
		found = true
		t.Name = task.Name
		tasks[index] = t
	}
	if !found {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func deleteTaskById(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	idStr := mux.Vars(request)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var found bool = false
	var updateTasks []Task
	for _, t := range tasks {
		if t.Id == id {
			found = true
			continue
		}
		updateTasks = append(updateTasks, t)
	}

	if !found {
		http.Error(writer, "Task not found", http.StatusNotFound)
		return
	}
	tasks = updateTasks
	writer.WriteHeader(http.StatusNoContent)
}

func createTask(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var task Task

	err := json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	task.Id = generateNewId()
	tasks = append(tasks, task)
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(task)
}

func getTaskById(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	idStr := mux.Vars(request)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		if task.Id == id {
			json.NewEncoder(writer).Encode(task)
			return
		}
	}
	http.Error(writer, "Task not found", http.StatusNotFound)
}

func getAllTasks(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(tasks)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
