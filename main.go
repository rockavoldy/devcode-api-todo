package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func ReturnJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "Application/json")

		next.ServeHTTP(rw, r)
	})
}

func main() {
	mysql_host := os.Getenv("MYSQL_HOST")
	if mysql_host == "" {
		mysql_host = "localhost"
	}
	mysql_user := os.Getenv("MYSQL_USER")
	if mysql_user == "" {
		mysql_user = "root"
	}
	mysql_password := os.Getenv("MYSQL_PASSWORD")
	if mysql_password == "" {
		mysql_password = "example"
	}
	mysql_dbname := os.Getenv("MYSQL_DBNAME")
	if mysql_dbname == "" {
		mysql_dbname = "akhmad_maulana_akbar"
	}

	db := ConnectDB(mysql_host, mysql_user, mysql_password, mysql_dbname)
	defer db.Close()
	repo := NewRepo(db)

	router := chi.NewRouter()

	router.Use(ReturnJson)

	router.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"mysql_host":   mysql_host,
			"mysql_user":   mysql_user,
			"mysql_pass":   mysql_password,
			"mysql_dbname": mysql_dbname,
		}
		dataJson, _ := json.Marshal(data)
		rw.Write([]byte(dataJson))
	})

	router.Get("/activity-groups", func(rw http.ResponseWriter, r *http.Request) {
		data, _ := repo.GetActivities()
		print := &PrintActivityGroups{
			Status:  "Success",
			Message: "Success",
			Data:    data,
		}

		resp, _ := json.Marshal(print)

		rw.WriteHeader(200)
		rw.Write([]byte(resp))
	})

	router.Get("/activity-groups/{activityId}", func(rw http.ResponseWriter, r *http.Request) {
		activityId := chi.URLParam(r, "activityId")
		print := &PrintActivtyGroup{}

		data, err := repo.GetActivity(activityId)
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
	})

	router.Post("/activity-groups", func(rw http.ResponseWriter, r *http.Request) {
		data := &ActivityGroup{}
		print := &PrintActivtyGroup{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Println(err)
		}

		if err := data.Validate(); err != nil {
			print.Status = "Bad Request"
			print.Message = err.Error()
			print.Data = map[string]interface{}{}
			rw.WriteHeader(400)
			resp, _ := json.Marshal(print)

			rw.Write([]byte(resp))
			return
		}

		dataInsert, _ := repo.InsertActivity(data)

		print.Status = "Success"
		print.Message = "Success"
		print.Data = dataInsert
		rw.WriteHeader(201)
		resp, _ := json.Marshal(print)

		rw.Write([]byte(resp))
	})

	router.Delete("/activity-groups/{activityId}", func(rw http.ResponseWriter, r *http.Request) {
		activityId := chi.URLParam(r, "activityId")
		print := &PrintActivtyGroup{}

		deleted, _ := repo.DeleteActivity(activityId)
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
	})

	router.Patch("/activity-groups/{activityId}", func(rw http.ResponseWriter, r *http.Request) {
		activityId := chi.URLParam(r, "activityId")
		print := &PrintActivtyGroup{}
		data := make(map[string]interface{})

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Println(err)
		}

		if _, ok := data["title"]; !ok {
			print.Status = "Bad Request"
			print.Message = ErrTitleNull.Error()
			print.Data = map[string]interface{}{}
			rw.WriteHeader(400)
			resp, _ := json.Marshal(print)

			rw.Write([]byte(resp))
			return
		}

		updatedData, err := repo.UpdateActivity(activityId, data)
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
	})

	log.Println("Listening on port :3030")
	if err := http.ListenAndServe(":3030", router); err != nil {
		log.Fatalln(err)
	}
}
