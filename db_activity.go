package main

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// Insert to activities
func (r *Repo) InsertActivity(activity *ActivityGroup) (int64, error) {
	sqlQuery, args, _ := sq.Insert("activities").Columns("email", "title").Values(activity.Email, activity.Title).ToSql()

	res, err := r.DB.Exec(sqlQuery, args...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Get activity
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

// Get activities
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

// Update activity
func (r *Repo) UpdateActivity(id interface{}, columns map[string]interface{}) (map[string]interface{}, error) {
	_, err := r.GetActivity(id)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}

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

// Delete activity
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
