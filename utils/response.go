package utils

type GeneralResponse struct {
	Action bool   `json:"action"`
	Error  string `json:"error,omitempty"`
	Info   string `json:"info,omitempty"`
}
