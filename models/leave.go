package models

import "time"

type LeaveType string

const(
		SickLeave LeaveType="sick"
		CasualLeave LeaveType="casual"
		Otherwise LeaveType="other"
)

type Leave struct {
	ID int `json:"id" db:"id"`
	UserID int `json:"-" db:"user_id"`
	LeaveFrom string `json:"leaveFrom" db:"leave_from"`
	LeaveTo string `json:"leaveTo" db:"leave_to"`
	LeaveType LeaveType `json:"leaveType" db:"leave_type"`
	Reason string `json:"reason" db:"reason"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

type LeaveStat struct {
	LeaveType LeaveType `json:"leave_type" db:"leave_type"`
	Taken     int    `json:"taken" db:"taken"`
	Allowed   int    `json:"allowed" db:"allowed"`
}