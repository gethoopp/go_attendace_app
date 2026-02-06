package modules

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Users struct {
	Id         int            `json:"id_users"`
	Rfid       int            `json:"rfid_id"`
	FirstName  string         `json:"id_first_name"`
	LastName   string         `json:"id_last_name"`
	Department string         `json:"id_departement"`
	Email      string         `json:"email_user"`
	Password   string         `json:"password_user"`
	DeleteResp gorm.DeletedAt `gorm:"index"`
	Role       string         `json:"role"`
}

type UserResponse struct {
	Id         int            `json:"id_users"`
	Rfid       int            `json:"rfid_id"`
	FirstName  string         `json:"id_first_name"`
	LastName   string         `json:"id_last_name"`
	Department string         `json:"id_departement"`
	Email      string         `json:"email_user"`
	DeleteResp gorm.DeletedAt `gorm:"index"`
	Role       string         `json:"role"`
}

type DataImage struct {
	ImageRef []byte `json:"image_hash"`
}

type LoginRequest struct {
	Email    string `json:"email_user"`
	Password string `json:"password_user"`
}

type ClaimsData struct {
	UserID   int    `json:"user_id"`
	NamaUser string `json:"name_user"`
	jwt.RegisteredClaims
}

type GetTotalWorker struct {
	TotalWorker int `json:"total_day_working"`
}

type GetTotalWorkerResponse struct {
	TotalWorker int `json:"total_day_working"`
}
