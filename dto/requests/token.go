package requests

type TokenReq struct {
	Duration int    `json:"duration"`
	Category string `json:"category"` // "hadir" atau "telat"
}

type SubmitToken struct {
	TokenCode string  `json:"token_code"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
