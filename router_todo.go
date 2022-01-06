package main

import (
	"devcode-api-todo/model"
	"devcode-api-todo/repo"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	repo  *repo.Repo
	wp    *workerPool
	mutex *sync.Mutex
}

func NewTodo(repo *repo.Repo, wp *workerPool, mtx *sync.Mutex) *Todo {
	return &Todo{
		repo:  repo,
		wp:    wp,
		mutex: mtx,
	}
}

func RouterTodo(router fiber.Router, repo *repo.Repo, wp *workerPool, mtx *sync.Mutex) {
	todo := NewTodo(repo, wp, mtx)
	router.Get("/", todo.list)
	router.Get("/:todoId", todo.get)
	router.Post("/", todo.create)
	router.Delete("/:todoId", todo.delete)
	router.Patch("/:todoId", todo.update)
}

func (t *Todo) list(c *fiber.Ctx) error {
	activityGroupId := c.Query("activity_group_id")
	data, _ := t.repo.GetTodos(activityGroupId)
	print := &model.PrintTodoItems{
		Status:  "Success",
		Message: "Success",
		Data:    data,
	}

	return c.JSON(print)
}

func (t *Todo) get(c *fiber.Ctx) error {
	todoId, _ := c.ParamsInt("todoId")
	print := &model.PrintTodoItem{}

	data, err := t.repo.GetTodo(todoId)
	if err != nil {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %d Not Found", todoId)
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(404)

		return c.JSON(print)
	}

	print.Status = "Success"
	print.Message = "Success"
	print.Data = data

	return c.JSON(print)
}

func (t *Todo) create(c *fiber.Ctx) error {
	var data model.TodoItem
	print := &model.PrintTodoItem{}

	c.BodyParser(&data)

	if data.ActivityGroupId == 0 {
		print.Status = "Bad Request"
		print.Message = model.ErrActivityGroupIdNull.Error()
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(400)

		return c.JSON(print)
	}

	if data.Title == "" {
		print.Status = "Bad Request"
		print.Message = model.ErrTitleNull.Error()
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(400)

		return c.JSON(print)
	}

	dataStruct := model.NewTodoItem(t.mutex, data.ActivityGroupId, data.Title, data.IsActive, model.PriorityStringToInt(data.Priority))
	t.wp.AddTask(func() {
		t.repo.InsertTodo(dataStruct)
	})

	print.Status = "Success"
	print.Message = "Success"
	print.Data = dataStruct.MapToInterface()

	c.Response().SetStatusCode(201)

	return c.JSON(print)
}

func (t *Todo) delete(c *fiber.Ctx) error {
	todoId, _ := c.ParamsInt("todoId")
	print := &model.PrintTodoItem{}
	deleted, _ := t.repo.DeleteTodo(todoId)

	if !deleted {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %d Not Found", todoId)
		c.Response().SetStatusCode(404)
		return c.JSON(print)
	}

	print.Status = "Success"
	print.Message = "Success"
	print.Data = map[string]interface{}{}

	return c.JSON(print)
}

func (t *Todo) update(c *fiber.Ctx) error {
	todoId, _ := c.ParamsInt("todoId")
	print := &model.PrintTodoItem{}
	data := make(map[string]interface{})

	c.BodyParser(&data)
	updatedData, err := t.repo.UpdateTodo(todoId, data)

	if err == model.ErrRecordNotFound {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %d Not Found", todoId)
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(404)

		return c.JSON(print)
	}

	print.Status = "Success"
	print.Message = "Success"
	print.Data = updatedData

	return c.JSON(print)
}
