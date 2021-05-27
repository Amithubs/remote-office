package dbHelper

import (
	"fmt"
	"github.com/tejashwikalptaru/remote-office/database"
	"github.com/tejashwikalptaru/remote-office/models"
	"time"
)

// InsertLeave save new leave in database
func  InsertLeave(userId int, leaveFrom time.Time ,leaveTo time.Time, reason string,leaveType string) error {
	SQL := `INSERT INTO user_leave(user_id, leave_from, leave_to, reason, leave_type) VALUES ($1, $2, $3, $4, $5)`
	_, err := database.RemoteOfficeDB.Exec(SQL, userId, leaveFrom, leaveTo, reason, leaveType)
	return err
}
func  ModifyLeave(userID int, leaveFrom time.Time ,leaveTo time.Time, reason string,leaveType string,leaveID int) error {
	SQL := `UPDATE user_leave SET leave_from=$1, leave_to=$2,reason=$3,leave_type=$4 WHERE id=$5 AND user_id=$6`
	_, err := database.RemoteOfficeDB.Exec(SQL, leaveFrom, leaveTo, reason, leaveType,leaveID,userID)
	return err
}
func  GetLeave(userID ,offset,limit int) ([]models.Leave, error) {
	SQL := `SELECT id,leave_from,leave_to,leave_type, created_at, reason from user_leave WHERE archived_at IS NULL AND user_id=$1  ORDER BY id DESC offset $2 limit $3`
	leaves:= make([]models.Leave,0)
	err := database.RemoteOfficeDB.Select(&leaves,SQL,userID,offset,limit)
	if err!=nil{
		return nil, err
}
	return leaves, nil
}
func GetLeaveStats(userID,year int)([]models.LeaveStat, error){
	SQL := `SELECT leave_type, count(*) as taken,allowed_leaves from user_leave as ul
         join allowed_leaves as al on ul.leave_type = al.leave_type where user_id = $1
  			AND ul.created_at >= $2
  			AND ul.created_at <= $3
  			AND archived_at IS NULL
  			AND al.year = $4 
			group by ul.leave_type, allowed_leaves
			order by ul.leave_type;`
	startDate := fmt.Sprintf("%d-01-01", year)
	endDate := fmt.Sprintf("%d-12-31", year)

	stat := make([]models.LeaveStat, 0)
	err := database.RemoteOfficeDB.Select(&stat, SQL, userID, startDate, endDate, year)
	return stat, err
}
