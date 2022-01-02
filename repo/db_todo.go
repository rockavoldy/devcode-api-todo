package repo

import (
	"database/sql"
	"devcode-api-todo/model"

	sq "github.com/Masterminds/squirrel"
)

// Insert todo
func (r *Repo) InsertTodo(todo map[string]interface{}) (map[string]interface{}, error) {
	res, err := r.DB.Exec(`INSERT INTO todos (activity_group_id, title) VALUES (?, ?)`, todo["activity_group_id"], todo["title"])
	if err != nil {
		return nil, err
	}

	lastInsertId, _ := res.LastInsertId()

	return r.GetTodo(lastInsertId)
}

// Get todo
func (r *Repo) GetTodo(id interface{}) (map[string]interface{}, error) {
	row := r.DB.QueryRowx(`SELECT * FROM todos WHERE id=?`, id)

	todoItem := &model.TodoItem{}
	err := row.StructScan(todoItem)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return todoItem.MapToInterface(), nil
}

// Get todos
func (r *Repo) GetTodos(query string) ([]map[string]interface{}, error) {
	rows, err := r.DB.Queryx(`SELECT * FROM todos WHERE activity_group_id=?`, query)
	if err != nil {
		return nil, err
	}

	todoItems := make([]map[string]interface{}, 0)
	for rows.Next() {
		todoItem := &model.TodoItem{}
		rows.StructScan(todoItem)

		todoItems = append(todoItems, todoItem.MapToInterface())
	}

	return todoItems, nil
}

// Update todo
func (r *Repo) UpdateTodo(id interface{}, columns map[string]interface{}) (map[string]interface{}, error) {
	_, err := r.GetTodo(id)
	if err == sql.ErrNoRows {
		return nil, model.ErrRecordNotFound
	}

	sqlQuery, args, _ := sq.Update("todos").Where(sq.Eq{"id": id}).SetMap(columns).ToSql()

	res, err := r.DB.Exec(sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	affected, _ := res.RowsAffected()
	if affected > 0 {
		return r.GetTodo(id)
	}

	return nil, model.ErrRecordNotFound
}

// Delete todo
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
