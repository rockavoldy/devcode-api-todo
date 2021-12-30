package main

import (
	"fmt"
	"log"

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
	activityId := c.Params("activityId")
	print := &PrintActivtyGroup{}

	data, err := a.repo.GetActivity(activityId)
	if err != nil {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %s Not Found", activityId)
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
	dataInsert, _ := a.repo.GetActivity(insertedId)

	print.Status = "Success"
	print.Message = "Success"
	print.Data = dataInsert
	c.Response().SetStatusCode(201)

	return c.JSON(print)
}

func (a *Activity) delete(c *fiber.Ctx) error {
	activityId := c.Params("activityId")
	print := &PrintActivtyGroup{}

	deleted, _ := a.repo.DeleteActivity(activityId)
	if !deleted {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %s Not Found", activityId)
		c.Response().SetStatusCode(404)
	} else {
		print.Status = "Success"
		print.Message = "Success"
	}
	print.Data = map[string]interface{}{}

	return c.JSON(print)
}

func (a *Activity) update(c *fiber.Ctx) error {
	activityId := c.Params("activityId")
	print := &PrintActivtyGroup{}
	data := make(map[string]interface{})

	c.BodyParser(&data)

	updatedData, err := a.repo.UpdateActivity(activityId, data)
	if err != nil {
		log.Println(err)
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %s Not Found", activityId)
		print.Data = map[string]interface{}{}
		c.Response().SetStatusCode(404)

		return c.JSON(print)
	}

	print.Status = "Success"
	print.Message = "Success"
	print.Data = updatedData

	return c.JSON(print)
}
