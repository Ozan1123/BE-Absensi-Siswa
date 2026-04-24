package responses

type ImportResult struct {
	Inserted     int      `json:"inserted"`
	Duplicates   int      `json:"duplicates"`
	Failed       int      `json:"failed"`
	SkippedUsers []string `json:"skipped_users"`
}
