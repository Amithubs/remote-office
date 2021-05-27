package dbHelper

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/tejashwikalptaru/remote-office/database"
	"github.com/tejashwikalptaru/remote-office/models"
	"time"
)

// IsUserExist checks if an user exists with a given email
func IsUserExist(email string) (int, error) {
	SQL := `SELECT id FROM users WHERE email = $1 AND archived_at IS NULL`
	var id int
	err := database.RemoteOfficeDB.Get(&id, SQL, email)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, nil
}

// InsertUser creates a new user entry in table
func InsertUser(name, email, password string, permission models.UserPermission) error {
	txError:=database.Tx(func(tx *sqlx.Tx) error {
		SQL := `INSERT INTO users(name, email, password) VALUES ($1, $2, $3)`
		var userID int
		err:=tx.Get(&userID,SQL,name, email,password)
		if err!=nil{
			return err
		}
		SQL = `INSERT INTO user_permission(user_id, permission_type) VALUES ($1, $2)`
		_,err=tx.Exec(SQL,userID,permission)
		if err!=nil{
			return err
			}
			return nil
	})
return txError
}

// GetUserPasswordByEmail returns the password for a given user with email
func GetUserPasswordByEmail(email string) (string, error) {
	SQL := `SELECT password FROM users WHERE email = $1 AND archived_at IS NULL`
	var password string
	err := database.RemoteOfficeDB.Get(&password, SQL, email)
	if err != nil {
		return "", err
	}
	return password, nil
}

// StoreUserSession saves a user session token in database
func StoreUserSession(userId int, token string) error {
	SQL := `INSERT INTO user_session(user_id, token, last_used_at) VALUES ($1, $2, $3)`
	_, err := database.RemoteOfficeDB.Exec(SQL, userId, token, time.Now())
	return err
}



// GetUserByToken gets the user details for a given token
func GetUserByToken(token string) (*models.User, error) {
	SQL := `SELECT 
				u.id,
				u.name, 
				u.phone, 
				u.email, 
				u.position,
				u.profile_image,
				u.created_at
			FROM users u 
			JOIN user_session us ON us.user_id = u.id
			WHERE u.archived_at IS NULL
			AND us.token = $1`
	var user models.User
	if err:= database.RemoteOfficeDB.Get(&user, SQL, token); err != nil {
		return nil, err
	}
	if user.ProfileImageID.Valid {
		//todo find user profile image, make public link and assign
		user.ProfileImageLink = "link"
	}
	return &user, nil
}
//DeleteUserSessionAll deletes all sessions for given userid
func DeleteUserSessionAll(userID int) error{
	SQL := `DELETE FROM user_session WHERE user_id=$1`
	_, err:= database.RemoteOfficeDB.Exec(SQL,userID)
	return err
}
func DeleteUserSession(token string,userID int) error{
	SQL := `DELETE FROM user_session WHERE user_id=$1 AND token=$2`
	_, err:= database.RemoteOfficeDB.Exec(SQL,userID,token)
	return err
}

func UserPermissionByID(userID int) ([]models.UserPermission,error){
	SQL:=`SELECT permission_type FROM user_permission WHERE user_id=$1`
	permissions:=make([]models.UserPermission,0)
	err:=database.RemoteOfficeDB.Select(&permissions,SQL,userID)
	if err!=nil{
		return nil, err
	}
	return permissions , nil
}
func UpdateUserInfo(userID int, phoneNo, position string) error {
	SQL := `UPDATE users set phone=$1,position=$2 WHERE id = $3`
	_, err := database.RemoteOfficeDB.Exec(SQL, phoneNo, position, userID)
	return err
}