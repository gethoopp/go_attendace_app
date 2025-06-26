package modules

type RequestChat struct {
	Prompt string `json:"prompt"`
}

type RequestStreamChat struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}
