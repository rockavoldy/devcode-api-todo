package main

import (
	"fmt"
	"log"

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
	todoId := c.Params("todoId")
	print := &PrintTodoItem{}

	data, err := t.repo.GetTodo(todoId)
	if err != nil {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %s Not Found", todoId)
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
	dataInsert, _ := t.repo.GetTodo(insertedId)

	print.Status = "Success"
	print.Message = "Success"
	print.Data = dataInsert

	c.Response().SetStatusCode(201)

	return c.JSON(print)
}

func (t *Todo) delete(c *fiber.Ctx) error {
	todoId := c.Params("todoId")
	print := &PrintTodoItem{}

	deleted, _ := t.repo.DeleteTodo(todoId)
	if !deleted {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %s Not Found", todoId)
		c.Response().SetStatusCode(404)
	} else {
		print.Status = "Success"
		print.Message = "Success"
	}

	print.Data = map[string]interface{}{}

	return c.JSON(print)
}

func (t *Todo) update(c *fiber.Ctx) error {
	todoId := c.Params("todoId")
	print := &PrintTodoItem{}
	data := make(map[string]interface{})

	c.BodyParser(&data)

	updatedData, err := t.repo.UpdateTodo(todoId, data)
	if err != nil {
		log.Println(err)
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Todo with ID %s Not Found", todoId)
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(404)

		return c.JSON(print)
	}

	print.Status = "Success"
	print.Message = "Success"
	print.Data = updatedData

	return c.JSON(print)
}
