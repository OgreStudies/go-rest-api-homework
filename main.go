package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Получить все задачи
func getTasks(res http.ResponseWriter, req *http.Request) {
	//Преобразование tasks в json
	resp, err := json.Marshal(&tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	//Формирование заголовка
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	//Запись тела
	res.Write(resp)
}

// Добавить задачу
func postTask(res http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer
	//Чтение тела
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	//Преобразование тела в Task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	//Добавление новой задачи
	tasks[task.ID] = task

	//Формирование заголовка
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
}

// Получить задачу по идентификатору
func getTask(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	//Пытаемся получить задачу по id
	task, ok := tasks[id]
	if !ok {
		http.Error(res, "Task not found", http.StatusBadRequest)
		return
	}

	//Преобразование task в json
	resp, err := json.Marshal(&task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	//Формирование заголовка
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	//Запись тела
	res.Write(resp)

}

// Удалить задачу
func deleteTask(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	//Пытаемся найти задачу по id
	_, ok := tasks[id]
	if !ok {
		http.Error(res, "Task not found", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	//Формирование заголовка
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
