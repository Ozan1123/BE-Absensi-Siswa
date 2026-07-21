package repo

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
)

func GetAttendanceRows(kelas, jurusan, startDate, endDate string) ([]models.Users, error) {

	start, end, err := utils.DateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var users []models.Users

	db := database.DB.Model(&models.Users{}).Where("role = ?", "siswa")

	if kelas != "" {
		db = db.Where("class_group = ?", kelas)
	}

	if jurusan != "" {
		db = db.Where("class_group LIKE ?", "%"+jurusan+"%")
	}

	err = db.Preload("AttedanceLogs", "clock_in_time >= ? AND clock_in_time <= ?", start, end).
		Order("full_name ASC").
		Find(&users).Error

	return users, err
}
