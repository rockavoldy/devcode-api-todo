package repo

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var createTableActivities = `
	CREATE TABLE IF NOT EXISTS activities (
		id INT PRIMARY KEY, 
		email VARCHAR(192), 
		title VARCHAR(192) NOT NULL, 
		created_at DATETIME, 
		updated_at DATETIME, 
		deleted_at DATETIME
	) ENGINE=InnoDB;`

var createTableTodos = `
	CREATE TABLE IF NOT EXISTS todos (
		id INT PRIMARY KEY, 
		activity_group_id INT NOT NULL, 
		title VARCHAR(192) NOT NULL, 
		is_active BOOLEAN DEFAULT TRUE, 
		priority ENUM('very-high', 'high', 'normal', 'low', 'very-low') DEFAULT 'very-high', 
		created_at DATETIME, 
		updated_at DATETIME, 
		deleted_at DATETIME
	) ENGINE=InnoDB;
`

func ConnectDB(host, user, pass, dbname string) *sqlx.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", user, pass, host, dbname)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatalln("Cannot connect to db server")
	}

	log.Println("Successfully connected to db server")
	db.MustExec(`DROP TABLE IF EXISTS todos`)
	db.MustExec(`DROP TABLE IF EXISTS activities`)

	db.MustExec(createTableActivities)
	db.MustExec(createTableTodos)

	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	return db
}

type Repo struct {
	DB *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		DB: db,
	}
}
