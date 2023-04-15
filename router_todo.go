package main

import (
	"devcode-api-todo/model"
	"devcode-api-todo/repo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Todo struct {
	repo *repo.Repo
}

func NewTodo(repo *repo.Repo) *Todo {
	return &Todo{
		repo: repo,
	}
}

func RouterTodo(repo *repo.Repo) http.Handler {
	todo := NewTodo(repo)
	router := chi.NewRouter()
	router.Get("/", todo.list)
	router.Get("/{todoId}", todo.get)
	router.Post("/", todo.create)
	router.Delete("/{todoId}", todo.delete)
	router.Patch("/{todoId}", todo.update)

	return router
}

func (t *Todo) list(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	activityGroupId := r.URL.Query().Get("activity_group_id")
	data, _ := t.repo.GetTodos(activityGroupId)
	print := &model.PrintTodoItems{
		Status:  "Success",
		Message: "Success",
		Data:    data,
	}

	resp, _ := json.Marshal(print)

	rw.WriteHeader(200)
	rw.Write([]byte(resp))
}

func (t *Todo) get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	todoParams := chi.URLParam(r, "todoId")
	todoId, _ := strconv.ParseInt(todoParams, 10, 64)
	print := &model.PrintTodoItem{}

	data, err := t.repo.GetTodo(todoId)
	if err != nil {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %d Not Found", todoId)
		print.Data = model.TodoItem{}
		rw.WriteHeader(404)

		resp, _ := json.Marshal(print)
		rw.Write([]byte(resp))
		return
	}

	print.Status = "Success"
	print.Message = "Success"
	print.Data = data
	rw.WriteHeader(200)
	resp, _ := json.Marshal(print)
	rw.Write([]byte(resp))
}

func (t *Todo) create(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	var data map[string]interface{}
	print := &model.PrintTodoItem{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	if _, ok := data["activity_group_id"]; !ok {
		print.Status = "Bad Request"
		print.Message = model.ErrActivityGroupIdNull.Error()
		print.Data = model.TodoItem{}
		rw.WriteHeader(400)
		resp, _ := json.Marshal(print)

		rw.Write([]byte(resp))
		return
	}

	if _, ok := data["title"]; !ok {
		print.Status = "Bad Request"
		print.Message = model.ErrTitleNull.Error()
		print.Data = model.TodoItem{}
		rw.WriteHeader(400)
		resp, _ := json.Marshal(print)

		rw.Write([]byte(resp))
		return
	}

	dataInsert, _ := t.repo.InsertTodo(data)

	print.Status = "Success"
	print.Message = "Success"
	print.Data = dataInsert
	rw.WriteHeader(201)
	resp, _ := json.Marshal(print)

	rw.Write([]byte(resp))
}

func (t *Todo) delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	todoParams := chi.URLParam(r, "todoId")
	todoId, _ := strconv.ParseInt(todoParams, 10, 64)
	print := &model.PrintTodoItem{}

	deleted, _ := t.repo.DeleteTodo(todoId)
	if !deleted {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %d Not Found", todoId)
		rw.WriteHeader(404)
	} else {
		print.Status = "Success"
		print.Message = "Success"
		rw.WriteHeader(200)
	}

	print.Data = model.TodoItem{}
	resp, _ := json.Marshal(print)
	rw.Write([]byte(resp))
}

func (t *Todo) update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	todoParams := chi.URLParam(r, "todoId")
	todoId, _ := strconv.ParseInt(todoParams, 10, 64)
	print := &model.PrintTodoItem{}
	data := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	updatedData, err := t.repo.UpdateTodo(todoId, data)
	if err != nil {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %d Not Found", todoId)
		print.Data = model.TodoItem{}
		rw.WriteHeader(404)
		resp, _ := json.Marshal(print)

		rw.Write([]byte(resp))
		return
	}

	print.Status = "Success"
	print.Message = "Success"
	rw.WriteHeader(200)
	print.Data = updatedData
	resp, _ := json.Marshal(print)
	rw.Write([]byte(resp))
}
