package responses

type UserRes struct {
	ID         int64  `json:"id"`
	Nisn       string `json:"nisn"`
	FullName   string `json:"full_name"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	ClassGroup string `json:"class_group"`
}
