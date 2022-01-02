package repo

import (
	"database/sql"
	"devcode-api-todo/model"

	sq "github.com/Masterminds/squirrel"
)

// Insert to activities
func (r *Repo) InsertActivity(activity map[string]string) (map[string]interface{}, error) {
	res, err := r.DB.Exec(`INSERT INTO activities (email, title) VALUES (?, ?)`, activity["email"], activity["title"])
	if err != nil {
		return nil, err
	}

	lastInsertId, _ := res.LastInsertId()

	return r.GetActivity(lastInsertId)
}

// Get activity
func (r *Repo) GetActivity(id interface{}) (map[string]interface{}, error) {
	row := r.DB.QueryRowx(`SELECT * FROM activities WHERE id=?`, id)

	activity := &model.ActivityGroup{}
	err := row.StructScan(activity)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return activity.MapToInterface(), nil
}

// Get activities
func (r *Repo) GetActivities() ([]map[string]interface{}, error) {
	rows, err := r.DB.Queryx(`SELECT * FROM activities`)
	if err != nil {
		return nil, err
	}

	activities := make([]map[string]interface{}, 0)
	for rows.Next() {
		activity := &model.ActivityGroup{}
		rows.StructScan(activity)

		activities = append(activities, activity.MapToInterface())
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

	res, err := r.DB.Exec(sqlQuery, args...)
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
