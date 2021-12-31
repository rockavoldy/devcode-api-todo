package main

import (
	"devcode-api-todo/model"
	"devcode-api-todo/repo"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	repo *repo.Repo
}

func NewTodo(repo *repo.Repo) *Todo {
	return &Todo{
		repo: repo,
	}
}

func RouterTodo(router fiber.Router, repo *repo.Repo) {
	todo := NewTodo(repo)
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
	data := make(map[string]interface{})
	print := &model.PrintTodoItem{}

	c.BodyParser(&data)

	if _, ok := data["activity_group_id"]; !ok {
		print.Status = "Bad Request"
		print.Message = model.ErrActivityGroupIdNull.Error()
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(400)

		return c.JSON(print)
	}

	if _, ok := data["title"]; !ok {
		print.Status = "Bad Request"
		print.Message = model.ErrTitleNull.Error()
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(400)

		return c.JSON(print)
	}

	dataInsert, _ := t.repo.InsertTodo(data)

	print.Status = "Success"
	print.Message = "Success"
	print.Data = dataInsert

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
