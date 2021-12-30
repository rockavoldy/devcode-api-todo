package main

import (
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
		mysql_dbname = "teestdb"
	}

	db := ConnectDB(mysql_host, mysql_user, mysql_password, mysql_dbname)
	repo := NewRepo(db)
	defer repo.DB.Close()

	router := chi.NewRouter()
	router.Use(ReturnJson)

	router.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello, World!\nDevcode Challenge #2"))
	})

	router.Mount("/activity-groups", RouterActivity(repo))
	router.Mount("/todo-items", RouterTodo(repo))

	log.Println("Listening on port :3030")
	if err := http.ListenAndServe(":3030", router); err != nil {
		log.Fatalln(err)
	}
}
