package main

import (
	"devcode-api-todo/repo"
	"log"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type WorkerPool interface {
	Run()
	AddTask(task func())
}

type workerPool struct {
	maxWorker  int
	queuedTask chan func()
}

func (wp *workerPool) Run() {
	for i := 0; i < wp.maxWorker; i++ {
		go func(workerID int) {
			for task := range wp.queuedTask {
				task()
			}
		}(i + 1)
	}
}

func (wp *workerPool) AddTask(task func()) {
	wp.queuedTask <- task
}

func NewWorkerPool(maxWorker int) *workerPool {
	return &workerPool{
		maxWorker:  maxWorker,
		queuedTask: make(chan func()),
	}
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

	db := repo.ConnectDB(mysql_host, mysql_user, mysql_password, mysql_dbname)
	repo := repo.NewRepo(db)
	defer repo.DB.Close()

	var mutex sync.Mutex
	wp := NewWorkerPool(100)
	wp.Run()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Prefork:               true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!\nDevcode challenge #2 with gofiber")
	})

	activity := app.Group("/activity-groups")
	RouterActivity(activity, repo, wp, &mutex)

	todo := app.Group("/todo-items")
	RouterTodo(todo, repo, wp, &mutex)

	log.Println("Listening on port :3030")

	if err := app.Listen(":3030"); err != nil {
		log.Fatalln(err)
	}
}
