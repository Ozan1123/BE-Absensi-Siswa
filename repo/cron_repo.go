package repo

import (
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"gorm.io/gorm"
)

type UnattendedStudent struct {
	ID          int64
	Nisn        string
	FullName    string
	ClassGroup  string
	ParentPhone string
}

type AbsentStudentWithStatus struct {
	ID          int64
	Nisn        string
	FullName    string
	ClassGroup  string
	ParentPhone string
	Status      string
}

// GetUnattendedStudents mengambil siswa yang belum absen sama sekali hari ini
func GetUnattendedStudents(db *gorm.DB) ([]UnattendedStudent, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
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

// GetStudentsByStatusToday mengambil siswa dengan status tertentu (sakit/izin) hari ini
func GetStudentsByStatusToday(db *gorm.DB, statuses []string) ([]AbsentStudentWithStatus, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour)

	var students []AbsentStudentWithStatus

	err := db.
		Table("users u").
		Select("u.id, u.nisn, u.full_name, u.class_group, u.parent_phone, l.status").
		Joins(`INNER JOIN (
			SELECT user_id, status
			FROM attedance_logs
			WHERE clock_in_time >= ? AND clock_in_time < ?
			AND status IN ?
		) l ON l.user_id = u.id`, start, end, statuses).
		Where("u.role = ?", "siswa").
		Order("u.class_group ASC, u.full_name ASC").
		Scan(&students).Error

	return students, err
}

// GetNotificationSettingsMap mengambil semua settings sebagai map
func GetNotificationSettingsMap(db *gorm.DB) (map[string]string, error) {
	var settings []models.NotificationSettings
	if err := db.Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.SettingKey] = s.SettingValue
	}
	return result, nil
}

// IsNotificationSentToday cek apakah notifikasi sudah dikirim hari ini untuk user tertentu
func IsNotificationSentToday(db *gorm.DB, userID int64, today string) bool {
	var count int64
	db.Model(&models.NotificationLogs{}).
		Where("user_id = ? AND sent_date = ?", userID, today).
		Count(&count)
	return count > 0
}

// TodayDateString mengembalikan tanggal hari ini format YYYY-MM-DD
func TodayDateString() string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc).Format("2006-01-02")
}
