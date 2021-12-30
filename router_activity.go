package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Activity struct {
	repo *Repo
}

func NewActivity(repo *Repo) *Activity {
	return &Activity{
		repo: repo,
	}
}

func RouterActivity(repo *Repo) http.Handler {
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
	print := &PrintActivityGroups{
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
	activityId := chi.URLParam(r, "activityId")
	print := &PrintActivtyGroup{}

	data, err := a.repo.GetActivity(activityId)
	if err != nil {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %s Not Found", activityId)
		print.Data = map[string]interface{}{}
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

func (a *Activity) create(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	data := &ActivityGroup{}
	print := &PrintActivtyGroup{}

	start := time.Now()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}
	log.Println("jsondecode: ", time.Since(start))

	start = time.Now()
	err = data.Validate()
	log.Println("validate: ", time.Since(start))
	if err != nil {
		print.Status = "Bad Request"
		print.Message = err.Error()
		print.Data = map[string]interface{}{}
		rw.WriteHeader(400)
		resp, _ := json.Marshal(print)

		rw.Write([]byte(resp))
		return
	}

	start = time.Now()
	insertedId, _ := a.repo.InsertActivity(data)
	dataInsert, _ := a.repo.GetActivity(insertedId)
	log.Println("insert-get: ", time.Since(start))

	print.Status = "Success"
	print.Message = "Success"
	print.Data = dataInsert
	rw.WriteHeader(201)
	resp, _ := json.Marshal(print)

	rw.Write([]byte(resp))
}

func (a *Activity) delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "Application/json")
	activityId := chi.URLParam(r, "activityId")
	print := &PrintActivtyGroup{}

	deleted, _ := a.repo.DeleteActivity(activityId)
	if !deleted {
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %s Not Found", activityId)
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
	activityId := chi.URLParam(r, "activityId")
	print := &PrintActivtyGroup{}
	data := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	updatedData, err := a.repo.UpdateActivity(activityId, data)
	if err != nil {
		log.Println(err)
		print.Status = "Not Found"
		print.Message = fmt.Sprintf("Activity with ID %s Not Found", activityId)
		print.Data = map[string]interface{}{}
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
