package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(host, user, pass, dbname string) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", user, pass, host, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatalln("Cannot connect to db server")
	}

	log.Println("Successfully connected to db server")
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	queryTableActivities := `CREATE TABLE IF NOT EXISTS activities(
		id INT AUTO_INCREMENT PRIMARY KEY, 
		email VARCHAR(255) NOT NULL, 
		title VARCHAR(255) NOT NULL, 
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, 
		deleted_at DATETIME) ENGINE=InnoDB;`

	queryTableTodos := `CREATE TABLE IF NOT EXISTS todos(
		id INT AUTO_INCREMENT PRIMARY KEY, 
		activity_group_id INT NOT NULL, 
		title VARCHAR(255) NOT NULL, 
		is_active BOOLEAN NOT NULL DEFAULT TRUE, 
		priority ENUM('very-high', 'high', 'normal', 'low', 'very-low') NOT NULL DEFAULT 'very-high', 
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, 
		deleted_at DATETIME, 
		FOREIGN KEY(activity_group_id) REFERENCES activities(id)) ENGINE=InnoDB;`

	db.ExecContext(context.Background(), queryTableActivities)
	db.ExecContext(context.Background(), queryTableTodos)

	return db
}

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		DB: db,
	}
}
