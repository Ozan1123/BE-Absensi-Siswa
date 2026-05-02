package models

//model users
type Users struct {
	ID          int64
	Nisn        string `gorm:"type:varchar(50);unique"`
	FullName    string
	Username    string `gorm:"type:varchar(50);unique"`
	Password    string
	Role        string
	ClassGroup  string
	ParentPhone string `gorm:"column:parent_phone;type:varchar(15)" json:"parent_phone"`

	AttedanceTokens []AttedanceTokens `gorm:"foreignKey:CreatedBy"`
	AttedanceLogs   []AttedanceLogs   `gorm:"foreignKey:UserID"`
}
