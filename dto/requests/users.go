package requests

type Login struct {
	Nisn     string `json:"nisn"`
	Password string `json:"password"`
}
