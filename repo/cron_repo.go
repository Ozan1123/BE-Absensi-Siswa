package repo

import (
	"time"

	"gorm.io/gorm"
)

type UnattendedStudent struct {
	ID          int64
	Nisn        string
	FullName    string
	ClassGroup  string
	ParentPhone string
}

func GetUnattendedStudents(db *gorm.DB) ([]UnattendedStudent, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	var students []UnattendedStudent

	err := db.
		Table("users u").
		Select("u.id, u.nisn, u.full_name, u.class_group, u.parent_phone").
		Joins(`LEFT JOIN (
			SELECT user_id, id
			FROM attedance_logs
			WHERE clock_in_time >= ? AND clock_in_time < ?
		) l ON l.user_id = u.id`, start, end).
		Where("u.role = ?", "siswa").
		Where("l.id IS NULL").
		Order("u.class_group ASC, u.full_name ASC").
		Scan(&students).Error

	return students, err
}
