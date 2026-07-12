package services

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/requests"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/responses"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
)

func GetLogService(UserID int64, query requests.LogQuery) (responses.LogResponse, error) {
	var logs []responses.LogResMini
	var total int64
	offset := (query.Page - 1) * query.Limit

	db := database.DB.Model(&models.AttedanceLogs{}).Where("user_id = ?", UserID)

	if query.Search != "" {
		db = db.Where("status LIKE ? ", "%"+query.Search+"%")
	}
	db.Count(&total)
	if err := db.Select("id", "status", "clock_in_time", "captured_ip").Order("clock_in_time DESC").Limit(query.Limit).Offset(offset).Find(&logs).Error; err != nil {
		return responses.LogResponse{}, err
	}

	totalPages := (total + int64(query.Limit) - 1) / int64(query.Limit)
	return responses.LogResponse{
		Message: "Found User Logs",
		Data: map[string]any{
			"logs" : logs,
			"total" : total,
			"page" : query.Page,
			"limit" : query.Limit,
			"total_pages" : totalPages,
		},
	}, nil
}