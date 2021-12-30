package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
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
	defer repo.DB.Close()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!\nDevcode challenge #2 with gofiber")
	})

	activity := app.Group("/activity-groups")
	RouterActivity(activity, repo)

	todo := app.Group("/todo-items")
	RouterTodo(todo, repo)

	log.Println("Listening on port :3030")

	if err := app.Listen(":3030"); err != nil {
		log.Fatalln(err)
	}
}
