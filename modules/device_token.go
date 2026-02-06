package modules

type DeviceTokenRequest struct {
	DeviceToken string `json:"device_token"`
	Platform    string `json:"platform"`
}
