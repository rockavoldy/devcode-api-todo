package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/go-chi/chi/v5"
)

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

	router := chi.NewRouter()

	router.Mount("/activity-groups", RouterActivity(repo))
	router.Mount("/todo-items", RouterTodo(repo))

	log.Println("Listening on port :3030")
	if err := http.ListenAndServe(":3030", router); err != nil {
		log.Fatalln(err)
	}
}
