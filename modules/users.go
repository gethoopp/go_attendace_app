package modules

type Userdata struct {
	Rfid   int    `json:"rfid_tag"`
	Name   string `json:"name_user"`
	Email  string `json:"email_user"`
	Divisi string `json:"divisi_user"`
}

type DataImage struct {
	ImageRef string `json:"image_path"`
}
