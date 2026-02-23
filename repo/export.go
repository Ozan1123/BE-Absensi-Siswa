package repo

import (
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
)

type ExportRow struct {
	Nisn        string
	FullName    string
	ClassGroup  string
	Status      *string
	ClockInTime *time.Time
}

func GetAttendanceRows(kelas, jurusan, tanggal string) ([]ExportRow, error) {

	start, end, err := utils.DayRange(tanggal)
	if err != nil {
		return nil, err
	}

	var rows []ExportRow

	db := database.DB.
		Table("users u").
		Select(`
			u.nisn,
			u.full_name,
			u.class_group,
			l.status,
			l.clock_in_time
		`).
		Joins(`
			LEFT JOIN (
				SELECT user_id, status, clock_in_time
				FROM attedance_logs
				WHERE clock_in_time >= ? AND clock_in_time < ?
			) l ON l.user_id = u.id
		`, start, end).
		Where("u.role = ?", "siswa")

	if kelas != "" {
		db = db.Where("u.class_group = ?", kelas)
	}

	if jurusan != "" {
		db = db.Where("u.class_group LIKE ?", "%"+jurusan+"%")
	}

	err = db.Order("u.full_name ASC").Scan(&rows).Error

	return rows, err
}
