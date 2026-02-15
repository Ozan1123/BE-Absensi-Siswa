package repo

import (
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/responses"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
)

func GetDashboardData() (*responses.DashboardResponse, error) {
	var (
		totalTokens      int64
		todayTokens      int64
		activeTokens     int64
		todayAttendances int64
	)

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	// total token
	if err := database.DB.Model(&models.AttedanceTokens{}).
		Count(&totalTokens).Error; err != nil {
		return nil, err
	}

	// token hari ini (FAST)
	if err := database.DB.Model(&models.AttedanceTokens{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&todayTokens).Error; err != nil {
		return nil, err
	}

	// active token
	if err := database.DB.Model(&models.AttedanceTokens{}).
		Where("is_active = ?", true).
		Count(&activeTokens).Error; err != nil {
		return nil, err
	}

	// attendance hari ini (FAST)
	if err := database.DB.Model(&models.AttedanceLogs{}).
		Where("clock_in_time >= ? AND clock_in_time < ?", start, end).
		Count(&todayAttendances).Error; err != nil {
		return nil, err
	}

	return &responses.DashboardResponse{
		TotalTokens:       int(totalTokens),
		TokenHariIni:      int(todayTokens),
		ActiveTokens:      int(activeTokens),
		TotalAbsenHariIni: int(todayAttendances),
	}, nil
}
