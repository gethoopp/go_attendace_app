package modules

import (
	"time"

	"google.golang.org/genproto/googleapis/type/datetime"
)

type Attendance struct {
	Id         int               `json:"id_attendance"`
	IdUser     int               `json:"id_users"`
	Check_in   time.Time         `json:"check_in"`
	Check_out  time.Time         `json:"check_out"`
	DateIn     datetime.DateTime `json:"attendance_at"`
	Status     string            `json:"status"`
	Created_at datetime.DateTime `json:"created_at"`
}
