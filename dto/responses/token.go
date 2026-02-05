package responses

import "time"

type TokenRes struct {
	ID         int64    `json:"id"`
	TokenCode  string   `json:"token_code"`
	CreatedBy  UserMini `json:"created_by"`
	IsActive   bool     `json:"is_active"`
	ValidUntil time.Time
	CreatedAt time.Time
}