package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	mysql_host := os.Getenv("MYSQL_HOST")
	if mysql_host == "" {
		mysql_host = "localhost"
	}
	mysql_user := os.Getenv("MYSQL_USER")
	if mysql_user == "" {
		mysql_user = "akhmad"
	}
	mysql_password := os.Getenv("MYSQL_PASSWORD")
	if mysql_password == "" {
		mysql_password = "akhmad"
	}
	mysql_dbname := os.Getenv("MYSQL_DBNAME")
	if mysql_dbname == "" {
		mysql_dbname = "akhmad_maulana_akbar"
	}

	router := http.NewServeMux()

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"mysql_host":   mysql_host,
			"mysql_user":   mysql_user,
			"mysql_pass":   mysql_password,
			"mysql_dbname": mysql_dbname,
		}
		dataJson, _ := json.Marshal(data)
		rw.Write([]byte(dataJson))
	})

	log.Println("Listening on port :8090")
	if err := http.ListenAndServe(":8090", router); err != nil {
		log.Fatalln(err)
	}

}
