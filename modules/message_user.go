package modules

type Message struct {
	Token string `json:"token" binding:"required"`
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

type SendMessageRequest struct {
	Message         *Message `json:"message,omitempty"`
	ValidateOnly    bool     `json:"validateOnly,omitempty"`
	ForceSendFields []string `json:"-"`
	NullFields      []string `json:"-"`
}
