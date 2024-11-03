package modules

type Userdata struct {
	Userid int    `json:"user_id"`
	Rfid   int    `json:"rfid_tag"`
	Name   string `json:"name_user"`
	Email  string `json:"email_user"`
	Divisi string `json:"divisi_user"`
	Tokens string `json:"token_user"`
}
