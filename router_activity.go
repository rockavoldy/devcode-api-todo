package main

import (
	"devcode-api-todo/model"
	"devcode-api-todo/repo"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Activity struct {
	repo *repo.Repo
}

func NewActivity(repo *repo.Repo) *Activity {
	return &Activity{
		repo: repo,
	}
}

func RouterActivity(router fiber.Router, repo *repo.Repo) {
	activity := NewActivity(repo)

	router.Get("/", activity.list)
	router.Get("/:activityId", activity.get)
	router.Post("/", activity.create)
	router.Delete("/:activityId", activity.delete)
	router.Patch("/:activityId", activity.update)
}

func (a *Activity) list(c *fiber.Ctx) error {
	data, _ := a.repo.GetActivities()
	print := &model.PrintActivityGroups{
		Status:  "Success",
		Message: "Success",
		Data:    data,
	}

	return c.JSON(print)
}

func (a *Activity) get(c *fiber.Ctx) error {
	activityId, _ := c.ParamsInt("activityId")
	print := &model.PrintActivtyGroup{}

	data, err := a.repo.GetActivity(activityId)

	if err != nil {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %d Not Found", activityId)
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(404)

		return c.JSON(print)
	}

	print.Status = "Success"
	print.Message = "Success"
	print.Data = data

	return c.JSON(print)
}

func (a *Activity) create(c *fiber.Ctx) error {
	data := make(map[string]string)
	print := &model.PrintActivtyGroup{}

	c.BodyParser(&data)

	if _, ok := data["title"]; !ok {
		print.Status = "Bad Request"
		print.Message = model.ErrTitleNull.Error()
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(400)

		return c.JSON(print)
	}

	dataInsert, _ := a.repo.InsertActivity(data)

	print.Status = "Success"
	print.Message = "Success"
	print.Data = dataInsert
	c.Response().SetStatusCode(201)

	return c.JSON(print)
}

func (a *Activity) delete(c *fiber.Ctx) error {
	activityId, _ := c.ParamsInt("activityId")
	print := &model.PrintActivtyGroup{}

	deleted, _ := a.repo.DeleteActivity(activityId)

	if !deleted {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %d Not Found", activityId)
		c.Response().SetStatusCode(404)
		return c.JSON(print)
	}

	print.Status = "Success"
	print.Message = "Success"
	print.Data = map[string]interface{}{}

	return c.JSON(print)
}

func (a *Activity) update(c *fiber.Ctx) error {
	activityId, _ := c.ParamsInt("activityId")
	print := &model.PrintActivtyGroup{}
	data := make(map[string]interface{})

	c.BodyParser(&data)

	updatedData, err := a.repo.UpdateActivity(activityId, data)

	if err == model.ErrRecordNotFound {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %d Not Found", activityId)
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(404)

		return c.JSON(print)
	}

	print.Status = "Success"
	print.Message = "Success"
	print.Data = updatedData

	return c.JSON(print)
}
