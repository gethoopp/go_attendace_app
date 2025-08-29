package modules

import "github.com/golang-jwt/jwt/v5"

type Users struct {
	Id         int    `json:"id_users"`
	Rfid       int    `json:"rfid_id"`
	FirstName  string `json:"id_first_name"`
	LastName   string `json:"id_last_name"`
	Department string `json:"id_departement"`
	Email      string `json:"email_user"`
	Password   string `json:"password_user"`
}

type DataImage struct {
	ImageRef string `json:"image_path"`
}

type ClaimsData struct {
	NamaUser string `json:"name_user"`
	jwt.RegisteredClaims
}
