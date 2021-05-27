package dbHelper

import (
	"github.com/tejashwikalptaru/remote-office/database"
	"github.com/tejashwikalptaru/remote-office/models"
	"time"
)

func  InsertStandUp(userId int, data string, date time.Time) error {
	SQL := `INSERT INTO standup(user_id, data, created_at) VALUES ($1, $2, $3)`
	_, err := database.RemoteOfficeDB.Exec(SQL, userId, data, date)
	return err
}
func  ModifyStandUp(userId int, data string, date time.Time,standUpID int) error {
	SQL := `UPDATE standup SET data=$1,created_at=$2 WHERE id=$3 AND user_id=$4`
	_, err := database.RemoteOfficeDB.Exec(SQL,data, date,standUpID,userId)
	return err
}
func  GetStandUp(userID ,offset,limit int) ([]models.StandUp, error) {
	SQL := `SELECT id, data, created_at from standup WHERE archived_at IS NULL AND user_id=$1 ORDER BY id offset $2 limit $3`
	standUps:= make([]models.StandUp,0)
	err := database.RemoteOfficeDB.Select(&standUps,SQL,userID,offset,limit)
	if err!=nil{
		return nil, err
	}
	return standUps, nil
}