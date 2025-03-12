package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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

// getTask возвращает все задачи, которые хранятся в мапе.
// Конечная точка /tasks.
// Метод GET.
// При успешном возвращает статус 200 OK.
// При ошибке возвращает статус 500 Internal Server Error.
func getTasks(resp http.ResponseWriter, req *http.Request) {

	jsonTasks, err := json.Marshal(tasks)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(jsonTasks)
}

// postTask принимает задачу в теле запроса и сохраняет ее в мапе.
// Конечная точка /tasks.
// Метод POST.
// При успешном запросе возвращает статус 201 Created.
// При ошибке возвращает статус 400 Bad Request.
func postTask(resp http.ResponseWriter, req *http.Request) {

	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
}

// getTask возвращает задачу с указанным в запросе пути ID, если такая есть в мапе.
// Если такого ID нет, возвращает соответствующий статус.
// Конечная точка /tasks/{id}.
// Метод GET.
// При успешном выполнении возвращает статус 200 OK.
// В случае ошибки или отсутствия задачи в мапе возвращает статус 400 Bad Request.
func getTask(resp http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(resp, "Task does not exist", http.StatusBadRequest)
		return
	}

	jsonTask, err := json.Marshal(task)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(jsonTask)
}

// deleteTask удаляет задачу из мапы по её ID, есть ли задача с таким ID есть в мапе.
// Если нет - возвращает соответствующий статус.
// Конечная точка /tasks/{id}.
// Метод DELETE.
// При успешном выполнении запроса возвращает статус 200 OK.
// В случае ошибки или отсутствия задачи в мапе возвращает  статус 400 Bad Request.
func deleteTask(resp http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(resp, "Task does not exist", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
}

func main() {

	// создаём новый роутер
	r := chi.NewRouter()

	// регистрируем в роутере эндпойнт `/tasks` с обработчиком `getTasks`
	r.Get("/tasks", getTasks)

	// регистрируем в роутере эндпойнт `/tasks` с обработчиком `postTask`
	r.Post("/tasks", postTask)

	// регистрируем в роутере эндпойнт `/tasks/{id}` с обработчиком `getTask`
	r.Get("/tasks/{id}", getTask)

	// регистрируем в роутере эндпойнт `/tasks/{id}` с обработчиком `deleteTask`
	r.Delete("/tasks/{id}", deleteTask)

	// запуск роутера
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
