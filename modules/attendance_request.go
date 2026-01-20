package modules

import (
	"time"
)

type Attendance struct {
	Id         int        `json:"id_attendance"`
	IdUser     int        `json:"user_id"`
	Check_in   *time.Time `json:"check_in"`
	Check_out  *time.Time `json:"check_out"`
	DateIn     *time.Time `json:"attendance_date"`
	Status     string     `json:"status"`
	Created_at *time.Time `json:"created_at"`
}

type AttendanceRequest struct {
	Id         int        `json:"id_attendance"`
	IdUser     int        `json:"user_id"`
	Check_in   *time.Time `json:"check_in"`
	Check_out  *time.Time `json:"check_out"`
	DateIn     string     `json:"attendance_date"`
	Status     string     `json:"status"`
	Created_at *time.Time `json:"created_at"`
}
