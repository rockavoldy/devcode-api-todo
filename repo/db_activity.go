package repo

import (
	"context"
	"database/sql"
	"devcode-api-todo/model"

	sq "github.com/Masterminds/squirrel"
)

// Insert to activities
func (r *Repo) InsertActivity(activity map[string]string) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Insert("activities").Columns("email", "title").Values(activity["email"], activity["title"]).ToSql()

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

	return r.GetActivity(lastInsertId)
}

// Get activity
func (r *Repo) GetActivity(id interface{}) (map[string]interface{}, error) {
	sqlQuery, args, _ := sq.Select("*").From("activities").Where(sq.Eq{"id": id}).ToSql()

	prep, err := r.DB.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer prep.Close()

	row := prep.QueryRowContext(context.Background(), args...)

	activity := &model.ActivityGroup{}
	err = row.Scan(&activity.ID, &activity.Email, &activity.Title, &activity.CreatedAt, &activity.UpdatedAt, &activity.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return activity.MapToInterface(), nil
}

// Get activities
func (r *Repo) GetActivities() ([]map[string]interface{}, error) {
	sqlQuery, _, _ := sq.Select("*").From("activities").ToSql()

	rows, err := r.DB.QueryContext(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := make([]map[string]interface{}, 0)
	for rows.Next() {
		activity := &model.ActivityGroup{}
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
		return nil, model.ErrRecordNotFound
	}

	sqlQuery, args, _ := sq.Update("activities").Where(sq.Eq{"id": id}).SetMap(columns).ToSql()

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
		return r.GetActivity(id)
	}

	return nil, model.ErrRecordNotFound
}

// Delete activity
func (r *Repo) DeleteActivity(id interface{}) (bool, error) {
	sqlQuery, args, _ := sq.Delete("activities").Where(sq.Eq{"id": id}).ToSql()

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
