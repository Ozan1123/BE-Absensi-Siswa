package services

import (
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/requests"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
)

// DetermineAttendanceStatus menentukan status absensi berdasarkan kondisi token.
// - Token expired → "telat"
// - Token valid tapi sudah melewati LateAfter → "telat"
// - Token valid dan belum melewati LateAfter → "hadir"
func DetermineAttendanceStatus(token *models.AttedanceTokens, isExpired bool) string {
	if isExpired {
		return "telat"
	}

	if utils.Now().After(token.LateAfter) {
		return "telat"
	}

	return "hadir"
}




func GetAttendanceTrend7Days() ([]requests.TrendResponse, error) {

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	var result []requests.TrendResponse

	for i := 6; i >= 0; i-- {

		day := now.AddDate(0, 0, -i)

		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
		end := start.Add(24 * time.Hour)

		var total int64

		// 🔥 PAKAI POLA YANG SUDAH TERBUKTI (seperti dashboard kamu)
		err := database.DB.Model(&models.AttedanceLogs{}).
			Where("clock_in_time >= ? AND clock_in_time < ?", start, end).
			Count(&total).Error

		if err != nil {
			return nil, err
		}

		result = append(result, requests.TrendResponse{
			Date:  start.Format("2006-01-02"),
			Total: int(total),
		})
	}

	return result, nil
}
