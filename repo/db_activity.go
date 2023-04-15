package repo

import (
	"devcode-api-todo/model"

	sq "github.com/Masterminds/squirrel"
)

// Insert to activities
func (r *Repo) InsertActivity(activity map[string]string) (model.ActivityGroup, error) {
	sqlQuery, args, _ := sq.Insert("activities").Columns("email", "title").Values(activity["email"], activity["title"]).ToSql()

	res := r.DB.MustExec(sqlQuery, args...)

	lastInsertId, _ := res.LastInsertId()

	return r.GetActivity(lastInsertId)
}

// Get activity
func (r *Repo) GetActivity(id int64) (model.ActivityGroup, error) {
	sqlQuery, args, _ := sq.Select("*").From("activities").Where(sq.Eq{"id": id}).ToSql()

	activity := model.ActivityGroup{}
	err := r.DB.Get(&activity, sqlQuery, args...)
	if err != nil {
		return model.ActivityGroup{}, err
	}

	return activity, nil
}

// Get activities
func (r *Repo) GetActivities() ([]model.ActivityGroup, error) {
	sqlQuery, _, _ := sq.Select("*").From("activities").ToSql()

	activities := []model.ActivityGroup{}
	err := r.DB.Select(&activities, sqlQuery)
	if err != nil {
		return nil, err
	}

	return activities, nil
}

// Update activity
func (r *Repo) UpdateActivity(id int64, columns map[string]interface{}) (model.ActivityGroup, error) {
	_, err := r.GetActivity(id)
	if err != nil {
		return model.ActivityGroup{}, model.ErrRecordNotFound
	}

	sqlQuery, args, _ := sq.Update("activities").Where(sq.Eq{"id": id}).SetMap(columns).ToSql()

	res := r.DB.MustExec(sqlQuery, args...)
	if err != nil {
		return model.ActivityGroup{}, err
	}

	affected, _ := res.RowsAffected()
	if affected > 0 {
		return r.GetActivity(id)
	}

	return model.ActivityGroup{}, model.ErrRecordNotFound
}

// Delete activity
func (r *Repo) DeleteActivity(id int64) (bool, error) {
	sqlQuery, args, _ := sq.Delete("activities").Where(sq.Eq{"id": id}).ToSql()

	res := r.DB.MustExec(sqlQuery, args...)

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return false, nil
	}

	return true, nil
}
