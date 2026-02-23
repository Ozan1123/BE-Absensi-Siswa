package models

//model users
type Users struct {
	ID         int64
	Nisn       string
	FullName   string
	Username   string
	Password   string
	Role       string
	ClassGroup string

	AttedanceTokens []AttedanceTokens `gorm:"foreignKey:CreatedBy"`
	AttedanceLogs   []AttedanceLogs   `gorm:"foreignKey:UserID"`
}
