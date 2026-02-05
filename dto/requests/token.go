package requests

type TokenReq struct {
	Duration int `json:"duration"`
}

type SubmitToken struct {
	TokenCode string `json:"token_code"`
}
