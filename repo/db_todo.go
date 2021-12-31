package repo

import (
	"context"
	"database/sql"
	"devcode-api-todo/model"

	sq "github.com/Masterminds/squirrel"
)

// Insert todo
func (r *Repo) InsertTodo(todo map[string]interface{}) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Insert("todos").Columns("activity_group_id", "title").Values(todo["activity_group_id"], todo["title"]).ToSql()

	prep, err := r.DB.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	res, err := prep.ExecContext(context.Background(), args...)
	if err != nil {
		return nil, err
	}

	lastInsertId, _ := res.LastInsertId()

	return r.GetTodo(lastInsertId)
}

// Get todo
func (r *Repo) GetTodo(id interface{}) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Select("*").From("todos").Where(sq.Eq{"id": id}).ToSql()

	prep, err := r.DB.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	row := prep.QueryRowContext(context.Background(), args...)

	todoItem := &model.TodoItem{}
	err = row.Scan(&todoItem.ID, &todoItem.ActivityGroupId, &todoItem.Title, &todoItem.IsActive, &todoItem.Priority, &todoItem.CreatedAt, &todoItem.UpdatedAt, &todoItem.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return todoItem.MapToInterface(), nil
}

// Get todos
func (r *Repo) GetTodos(query string) ([]map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Select("*").From("todos").Where(sq.Eq{"activity_group_id": query}).ToSql()

	prep, err := r.DB.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	rows, err := prep.QueryContext(context.Background(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todoItems := make([]map[string]interface{}, 0)
	for rows.Next() {
		todoItem := &model.TodoItem{}
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
		return nil, model.ErrRecordNotFound
	}

	sqlQuery, args, _ := sq.Update("todos").Where(sq.Eq{"id": id}).SetMap(columns).ToSql()

	prep, err := r.DB.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	res, err := prep.ExecContext(context.Background(), args...)
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

	prep, err := r.DB.Prepare(sqlQuery)
	if err != nil {
		return false, err
	}
	defer prep.Close()

	res, err := prep.ExecContext(context.Background(), args...)
	if err != nil {
		return false, err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return false, nil
	}

	return true, nil
}
