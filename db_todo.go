package main

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// Insert todo
func (r *Repo) InsertTodo(todo map[string]interface{}) (int64, error) {
	sqlQuery, args, _ := sq.Insert("todos").Columns("activity_group_id", "title").Values(todo["activity_group_id"], todo["title"]).ToSql()

	conn, err := r.DB.Conn(context.Background())
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	res, err := conn.ExecContext(context.Background(), sqlQuery, args...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Get todo
func (r *Repo) GetTodo(id interface{}) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Select("*").From("todos").Where(sq.Eq{"id": id}).ToSql()

	conn, err := r.DB.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	row := conn.QueryRowContext(context.Background(), sqlQuery, args...)

	todoItem := &TodoItem{}
	err = row.Scan(&todoItem.ID, &todoItem.ActivityGroupId, &todoItem.Title, &todoItem.IsActive, &todoItem.Priority, &todoItem.CreatedAt, &todoItem.UpdatedAt, &todoItem.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return todoItem.MapToInterface(), nil
}

// Get todos
func (r *Repo) GetTodos(query string) ([]map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Select("*").From("todos").Where(sq.Eq{"activity_group_id": query}).ToSql()

	conn, err := r.DB.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(context.Background(), sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todoItems := make([]map[string]interface{}, 0)
	for rows.Next() {
		todoItem := &TodoItem{}
		rows.Scan(&todoItem.ID, &todoItem.ActivityGroupId, &todoItem.Title, &todoItem.IsActive, &todoItem.Priority, &todoItem.CreatedAt, &todoItem.UpdatedAt, &todoItem.DeletedAt)

		todoItemMap := todoItem.MapToInterface()

		todoItems = append(todoItems, todoItemMap)
	}

	return todoItems, nil
}

// Update todo
func (r *Repo) UpdateTodo(id interface{}, columns map[string]interface{}) (map[string]interface{}, error) {
	_, err := r.GetTodo(id)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}

	conn, err := r.DB.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sqlQuery, args, _ := sq.Update("todos").Where(sq.Eq{"id": id}).SetMap(columns).ToSql()

	res, err := conn.ExecContext(context.Background(), sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	affected, _ := res.RowsAffected()
	if affected > 0 {
		return r.GetTodo(id)
	}

	return nil, ErrRecordNotFound
}

// Delete todo
func (r *Repo) DeleteTodo(id interface{}) (bool, error) {
	sqlQuery, args, _ := sq.Delete("todos").Where(sq.Eq{"id": id}).ToSql()

	conn, err := r.DB.Conn(context.Background())
	if err != nil {
		return false, err
	}
	defer conn.Close()

	res, err := conn.ExecContext(context.Background(), sqlQuery, args...)
	if err != nil {
		return false, err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return false, nil
	}

	return true, nil
}
