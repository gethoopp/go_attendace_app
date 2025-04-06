package modules

import "github.com/golang-jwt/jwt/v5"

type Userdata struct {
	Rfid   int    `json:"rfid_tag"`
	Name   string `json:"name_user"`
	Email  string `json:"email_user"`
	Divisi string `json:"divisi_user"`
}

type RegUserdata struct {
	Rfid     int    `json:"rfid_id"`
	Name     string `json:"id_first_name"`
	LastName string `json:"id_last_name"`
	Divisi   string `json:"id_departement"`
	Email    string `json:"email_user"`
	Password string `json:"password_user"`
}

type DataImage struct {
	ImageRef string `json:"image_path"`
}

type ClaimsData struct {
	NamaUser string `json:"name_user"`
	jwt.RegisteredClaims
}
