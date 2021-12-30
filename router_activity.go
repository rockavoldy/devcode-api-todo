package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Activity struct {
	repo *Repo
}

func NewActivity(repo *Repo) *Activity {
	return &Activity{
		repo: repo,
	}
}

func RouterActivity(router fiber.Router, repo *Repo) {
	activity := NewActivity(repo)

	router.Get("/", activity.list)
	router.Get("/:activityId", activity.get)
	router.Post("/", activity.create)
	router.Delete("/:activityId", activity.delete)
	router.Patch("/:activityId", activity.update)
}

func (a *Activity) list(c *fiber.Ctx) error {
	data, _ := a.repo.GetActivities()
	print := &PrintActivityGroups{
		Status:  "Success",
		Message: "Success",
		Data:    data,
	}

	return c.JSON(print)
}

func (a *Activity) get(c *fiber.Ctx) error {
	activityId, _ := c.ParamsInt("activityId")
	print := &PrintActivtyGroup{}

	if !a.repo.inmem[activityId] {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %d Not Found", activityId)
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(404)

		return c.JSON(print)
	}

	data, _ := a.repo.GetActivity(activityId)

	print.Status = "Success"
	print.Message = "Success"
	print.Data = data

	return c.JSON(print)
}

func (a *Activity) create(c *fiber.Ctx) error {
	data := new(ActivityGroup)
	print := &PrintActivtyGroup{}

	c.BodyParser(&data)

	if err := data.Validate(); err != nil {
		print.Status = "Bad Request"
		print.Message = err.Error()
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(400)

		return c.JSON(print)
	}

	insertedId, _ := a.repo.InsertActivity(data)
	a.repo.Add(int(insertedId))
	dataInsert, _ := a.repo.GetActivity(insertedId)

	print.Status = "Success"
	print.Message = "Success"
	print.Data = dataInsert
	c.Response().SetStatusCode(201)

	return c.JSON(print)
}

func (a *Activity) delete(c *fiber.Ctx) error {
	activityId, _ := c.ParamsInt("activityId")
	print := &PrintActivtyGroup{}

	if !a.repo.inmem[activityId] {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %d Not Found", activityId)
		c.Response().SetStatusCode(404)
		return c.JSON(print)
	}

	a.repo.DeleteActivity(activityId)
	print.Status = "Success"
	print.Message = "Success"
	print.Data = map[string]interface{}{}

	return c.JSON(print)
}

func (a *Activity) update(c *fiber.Ctx) error {
	activityId, _ := c.ParamsInt("activityId")
	print := &PrintActivtyGroup{}
	data := make(map[string]interface{})

	c.BodyParser(&data)

	if !a.repo.inmem[activityId] {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %d Not Found", activityId)
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(404)

		return c.JSON(print)
	}

	updatedData, _ := a.repo.UpdateActivity(activityId, data)
	print.Status = "Success"
	print.Message = "Success"
	print.Data = updatedData

	return c.JSON(print)
}
