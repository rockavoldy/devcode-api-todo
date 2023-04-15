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

type Activity struct {
	repo *repo.Repo
}

func NewActivity(repo *repo.Repo) *Activity {
	return &Activity{
		repo: repo,
	}
}

func RouterActivity(repo *repo.Repo) http.Handler {
	activity := NewActivity(repo)

	router := chi.NewRouter()
	router.Get("/", activity.list)
	router.Get("/{activityId}", activity.get)
	router.Post("/", activity.create)
	router.Delete("/{activityId}", activity.delete)
	router.Patch("/{activityId}", activity.update)

	return router
}

func (a *Activity) list(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	data, _ := a.repo.GetActivities()
	print := &model.PrintActivityGroups{
		Status:  "Success",
		Message: "Success",
		Data:    data,
	}

	resp, _ := json.Marshal(print)

	rw.WriteHeader(200)
	rw.Write([]byte(resp))
}

func (a *Activity) get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	activityParams := chi.URLParam(r, "activityId")
	activityId, _ := strconv.ParseInt(activityParams, 10, 64)
	print := &model.PrintActivityGroup{}

	data, err := a.repo.GetActivity(activityId)
	if err != nil {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %d Not Found", activityId)
		print.Data = map[string]interface{}{}
		rw.WriteHeader(404)

		resp, _ := json.Marshal(print)
		rw.Write([]byte(resp))
		return
	}

	print.Status = "Success"
	print.Message = "Success"
	print.Data = data.MapToInterface()
	rw.WriteHeader(200)
	resp, _ := json.Marshal(print)
	rw.Write([]byte(resp))
}

func (a *Activity) create(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	var data map[string]string
	print := &model.PrintActivityGroup{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	if _, titleOk := data["title"]; !titleOk {
		print.Status = "Bad Request"
		print.Message = model.ErrTitleNull.Error()
		print.Data = map[string]interface{}{}
		rw.WriteHeader(400)
		resp, _ := json.Marshal(print)

		rw.Write([]byte(resp))
		return
	}

	dataInsert, _ := a.repo.InsertActivity(data)

	print.Status = "Success"
	print.Message = "Success"
	print.Data = dataInsert.MapToInterface()
	rw.WriteHeader(201)
	resp, _ := json.Marshal(print)

	rw.Write([]byte(resp))
}

func (a *Activity) delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	activityParams := chi.URLParam(r, "activityId")
	activityId, _ := strconv.ParseInt(activityParams, 10, 64)
	print := &model.PrintActivityGroup{}

	deleted, _ := a.repo.DeleteActivity(activityId)
	if !deleted {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %d Not Found", activityId)
		rw.WriteHeader(404)
	} else {
		print.Status = "Success"
		print.Message = "Success"
		rw.WriteHeader(200)
	}

	print.Data = map[string]interface{}{}
	resp, _ := json.Marshal(print)
	rw.Write([]byte(resp))
}

func (a *Activity) update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	activityParams := chi.URLParam(r, "activityId")
	activityId, _ := strconv.ParseInt(activityParams, 10, 64)
	print := &model.PrintActivityGroup{}
	var data map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	updatedData, err := a.repo.UpdateActivity(activityId, data)
	if err != nil {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %d Not Found", activityId)
		print.Data = map[string]interface{}{}
		rw.WriteHeader(404)
		resp, _ := json.Marshal(print)

		rw.Write([]byte(resp))
		return
	}

	print.Status = "Success"
	print.Message = "Success"
	rw.WriteHeader(200)
	print.Data = updatedData.MapToInterface()
	resp, _ := json.Marshal(print)
	rw.Write([]byte(resp))
}
