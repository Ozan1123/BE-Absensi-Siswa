package requests

type TrendRow struct {
	Date  string `json:"date"`
	Total int   `json:"total"`
}

type TrendResponse struct {
	Date  string `json:"date"`
	Total int    `json:"total"`
}
