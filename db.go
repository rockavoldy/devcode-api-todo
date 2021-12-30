package main

import (
	"database/sql"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(host, user, pass, dbname string) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Println(err)
	}

	log.Println("MySQL Version: ", version)

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

// Activities Table
func (r *Repo) InsertActivity(activity *ActivityGroup) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Insert("activities").Columns("email", "title").Values(activity.Email, activity.Title).ToSql()

	res, err := r.DB.Exec(sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	lastId, _ := res.LastInsertId()

	return r.GetActivity(int(lastId))
}

func (r *Repo) GetActivity(id interface{}) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Select("*").From("activities").Where(sq.Eq{"id": id}).ToSql()

	row := r.DB.QueryRow(sqlQuery, args...)

	var activity ActivityGroup
	err := row.Scan(&activity.ID, &activity.Email, &activity.Title, &activity.CreatedAt, &activity.UpdatedAt, &activity.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return activity.MapToInterface(), nil
}

func (r *Repo) GetActivities() ([]map[string]interface{}, error) {
	sqlQuery, _, _ := sq.Select("*").From("activities").ToSql()

	prep, err := r.DB.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}

	rows, err := prep.Query()
	if err != nil {
		return nil, err
	}

	activities := make([]map[string]interface{}, 0)
	for rows.Next() {
		var activity ActivityGroup
		rows.Scan(&activity.ID, &activity.Email, &activity.Title, &activity.CreatedAt, &activity.UpdatedAt, &activity.DeletedAt)

		activityMap := activity.MapToInterface()

		activities = append(activities, activityMap)
	}

	return activities, nil
}

func (r *Repo) UpdateActivity(id interface{}, columns map[string]interface{}) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Update("activities").Where(sq.Eq{"id": id}).SetMap(columns).ToSql()

	res, err := r.DB.Exec(sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	affected, _ := res.RowsAffected()
	if affected > 0 {
		return r.GetActivity(id)
	}

	return nil, ErrRecordNotFound
}

func (r *Repo) DeleteActivity(id interface{}) (bool, error) {
	sqlQuery, args, _ := sq.Delete("activities").Where(sq.Eq{"id": id}).ToSql()

	res, err := r.DB.Exec(sqlQuery, args...)
	if err != nil {
		return false, err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return false, nil
	}

	return true, nil
}

// Todos Table
func (r *Repo) InsertTodo(todo *TodoItem) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Insert("todos").Columns("activity_group_id", "title").Values(todo.ActivityGroupId, todo.Title).ToSql()

	res, err := r.DB.Exec(sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	lastId, _ := res.LastInsertId()

	return r.GetTodo(int(lastId))
}

func (r *Repo) GetTodo(id interface{}) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Select("*").From("todos").Where(sq.Eq{"id": id}).ToSql()

	row := r.DB.QueryRow(sqlQuery, args...)

	var todoItem TodoItem
	err := row.Scan(&todoItem.ID, &todoItem.ActivityGroupId, &todoItem.Title, &todoItem.IsActive, &todoItem.Priority, &todoItem.CreatedAt, &todoItem.UpdatedAt, &todoItem.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return todoItem.MapToInterface(), nil
}

func (r *Repo) GetTodos(query string) ([]map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Select("*").From("todos").Where(sq.Eq{"activity_group_id": query}).ToSql()

	prep, err := r.DB.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}

	rows, err := prep.Query(args...)
	if err != nil {
		return nil, err
	}

	todoItems := make([]map[string]interface{}, 0)
	for rows.Next() {
		var todoItem TodoItem
		rows.Scan(&todoItem.ID, &todoItem.ActivityGroupId, &todoItem.Title, &todoItem.IsActive, &todoItem.Priority, &todoItem.CreatedAt, &todoItem.UpdatedAt, &todoItem.DeletedAt)

		todoItemMap := todoItem.MapToInterface()

		todoItems = append(todoItems, todoItemMap)
	}

	return todoItems, nil
}

func (r *Repo) UpdateTodo(id interface{}, columns map[string]interface{}) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Update("todos").Where(sq.Eq{"id": id}).SetMap(columns).ToSql()

	res, err := r.DB.Exec(sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	affected, _ := res.RowsAffected()
	if affected > 0 {
		return r.GetTodo(id)
	}

	return nil, ErrRecordNotFound
}

func (r *Repo) DeleteTodo(id interface{}) (bool, error) {
	sqlQuery, args, _ := sq.Delete("todos").Where(sq.Eq{"id": id}).ToSql()

	res, err := r.DB.Exec(sqlQuery, args...)
	if err != nil {
		return false, err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return false, nil
	}

	return true, nil
}
