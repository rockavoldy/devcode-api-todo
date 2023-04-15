package repo

import (
	"devcode-api-todo/model"

	sq "github.com/Masterminds/squirrel"
)

// Insert todo
func (r *Repo) InsertTodo(todo map[string]interface{}) (model.TodoItem, error) {
	sqlQuery, args, _ := sq.Insert("todos").Columns("activity_group_id", "title").Values(todo["activity_group_id"], todo["title"]).ToSql()

	res := r.DB.MustExec(sqlQuery, args...)

	lastInsertId, _ := res.LastInsertId()

	return r.GetTodo(lastInsertId)
}

// Get todo
func (r *Repo) GetTodo(id int64) (model.TodoItem, error) {
	sqlQuery, args, _ := sq.Select("*").From("todos").Where(sq.Eq{"id": id}).ToSql()

	todoItem := model.TodoItem{}
	err := r.DB.Get(&todoItem, sqlQuery, args...)

	if err != nil {
		return model.TodoItem{}, err
	}

	return todoItem, nil
}

// Get todos
func (r *Repo) GetTodos(query string) ([]model.TodoItem, error) {
	sqlQuery, args, _ := sq.Select("*").From("todos").Where(sq.Eq{"activity_group_id": query}).ToSql()

	todoItems := []model.TodoItem{}
	err := r.DB.Select(&todoItems, sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	return todoItems, nil
}

// Update todo
func (r *Repo) UpdateTodo(id int64, columns map[string]interface{}) (model.TodoItem, error) {
	_, err := r.GetTodo(id)
	if err != nil {
		return model.TodoItem{}, model.ErrRecordNotFound
	}

	sqlQuery, args, _ := sq.Update("todos").Where(sq.Eq{"id": id}).SetMap(columns).ToSql()

	res := r.DB.MustExec(sqlQuery, args...)

	affected, _ := res.RowsAffected()
	if affected > 0 {
		return r.GetTodo(id)
	}

	return model.TodoItem{}, model.ErrRecordNotFound
}

// Delete todo
func (r *Repo) DeleteTodo(id int64) (bool, error) {
	sqlQuery, args, _ := sq.Delete("todos").Where(sq.Eq{"id": id}).ToSql()

	res := r.DB.MustExec(sqlQuery, args...)

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return false, nil
	}

	return true, nil
}
