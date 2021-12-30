package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	repo *Repo
}

func NewTodo(repo *Repo) *Todo {
	return &Todo{
		repo: repo,
	}
}

func RouterTodo(router fiber.Router, repo *Repo) {
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
	print := &PrintTodoItems{
		Status:  "Success",
		Message: "Success",
		Data:    data,
	}

	return c.JSON(print)
}

func (t *Todo) get(c *fiber.Ctx) error {
	todoId, _ := c.ParamsInt("todoId")
	print := &PrintTodoItem{}

	if !t.repo.inmem[todoId] {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %d Not Found", todoId)
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(404)

		return c.JSON(print)
	}

	data, _ := t.repo.GetTodo(todoId)
	print.Status = "Success"
	print.Message = "Success"
	print.Data = data

	return c.JSON(print)
}

func (t *Todo) create(c *fiber.Ctx) error {
	data := new(TodoItem)
	print := &PrintTodoItem{}

	c.BodyParser(&data)

	if err := data.Validate(); err != nil {
		print.Status = "Bad Request"
		print.Message = err.Error()
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(400)

		return c.JSON(print)
	}

	insertedId, _ := t.repo.InsertTodo(data)
	t.repo.Add(int(insertedId))
	dataInsert, _ := t.repo.GetTodo(insertedId)

	print.Status = "Success"
	print.Message = "Success"
	print.Data = dataInsert

	c.Response().SetStatusCode(201)

	return c.JSON(print)
}

func (t *Todo) delete(c *fiber.Ctx) error {
	todoId, _ := c.ParamsInt("todoId")
	print := &PrintTodoItem{}

	if !t.repo.inmem[todoId] {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %d Not Found", todoId)
		c.Response().SetStatusCode(404)
		return c.JSON(print)
	}

	t.repo.DeleteTodo(todoId)
	print.Status = "Success"
	print.Message = "Success"
	print.Data = map[string]interface{}{}

	return c.JSON(print)
}

func (t *Todo) update(c *fiber.Ctx) error {
	todoId, _ := c.ParamsInt("todoId")
	print := &PrintTodoItem{}
	data := make(map[string]interface{})

	c.BodyParser(&data)

	if !t.repo.inmem[todoId] {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %d Not Found", todoId)
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(404)

		return c.JSON(print)
	}

	updatedData, _ := t.repo.UpdateTodo(todoId, data)
	print.Status = "Success"
	print.Message = "Success"
	print.Data = updatedData

	return c.JSON(print)
}
