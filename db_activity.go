package main

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// Insert to activities
func (r *Repo) InsertActivity(activity map[string]string) (int64, error) {
	sqlQuery, args, _ := sq.Insert("activities").Columns("email", "title").Values(activity["email"], activity["title"]).ToSql()

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

// Get activity
func (r *Repo) GetActivity(id interface{}) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Select("*").From("activities").Where(sq.Eq{"id": id}).ToSql()

	conn, err := r.DB.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	row := conn.QueryRowContext(context.Background(), sqlQuery, args...)

	var activity map[string]interface{}
	err = row.Scan(activity["id"], activity["email"], activity["title"], activity["created_at"], activity["updated_at"], activity["deleted_at"])

	if err == sql.ErrNoRows {
		return nil, err
	}

	return activity, nil
}

// Get activities
func (r *Repo) GetActivities() ([]map[string]interface{}, error) {
	sqlQuery, _, _ := sq.Select("*").From("activities").ToSql()

	conn, err := r.DB.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	prep, err := conn.PrepareContext(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	rows, err := prep.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []map[string]interface{}
	for rows.Next() {
		var activity map[string]interface{}
		rows.Scan(activity["id"], activity["email"], activity["title"], activity["created_at"], activity["updated_at"], activity["deleted_at"])

		activityMap := activity

		activities = append(activities, activityMap)
	}

	return activities, nil
}

// Update activity
func (r *Repo) UpdateActivity(id interface{}, columns map[string]interface{}) (map[string]interface{}, error) {
	_, err := r.GetActivity(id)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}

	sqlQuery, args, _ := sq.Update("activities").Where(sq.Eq{"id": id}).SetMap(columns).ToSql()

	conn, err := r.DB.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res, err := conn.ExecContext(context.Background(), sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	affected, _ := res.RowsAffected()
	if affected > 0 {
		return r.GetActivity(id)
	}

	return nil, ErrRecordNotFound
}

// Delete activity
func (r *Repo) DeleteActivity(id interface{}) (bool, error) {
	sqlQuery, args, _ := sq.Delete("activities").Where(sq.Eq{"id": id}).ToSql()

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
