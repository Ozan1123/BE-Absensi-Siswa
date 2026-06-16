package handlers

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/responses"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/mappers"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/gofiber/fiber/v2"
)

// GetAttendanceLogsAdmin godoc
// @Summary Ambil riwayat absensi (Admin/Guru)
// @Description Mengambil riwayat absensi siswa secara paginated dengan filter rentang tanggal, kelas, status, dan pencarian nama/NISN
// @Tags attendance
// @Produce json
// @Param page query int false "Halaman (default: 1)"
// @Param limit query int false "Batas item per halaman (default: 20)"
// @Param start_date query string false "Filter tanggal mulai (Format: YYYY-MM-DD)"
// @Param end_date query string false "Filter tanggal selesai (Format: YYYY-MM-DD)"
// @Param class_group query string false "Filter berdasarkan Kelas"
// @Param status query string false "Filter berdasarkan Status (hadir, telat, alfa, sakit, izin)"
// @Param search query string false "Pencarian berdasarkan Nama atau NISN Siswa"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /attendance/logs [get]
func GetAttendanceLogsAdmin(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	classGroup := c.Query("class_group")
	status := c.Query("status")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	query := database.DB.Model(&models.AttedanceLogs{}).
		Joins("JOIN users ON users.id = attedance_logs.user_id").
		Select("attedance_logs.*")

	if startDate != "" {
		query = query.Where("attedance_logs.clock_in_time >= ?", startDate+" 00:00:00")
	}
	if endDate != "" {
		query = query.Where("attedance_logs.clock_in_time <= ?", endDate+" 23:59:59")
	}
	if classGroup != "" {
		query = query.Where("users.class_group = ?", classGroup)
	}
	if status != "" {
		query = query.Where("attedance_logs.status = ?", status)
	}
	if search != "" {
		query = query.Where("users.full_name LIKE ? OR users.nisn LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal menghitung total log absensi"})
	}

	var logs []models.AttedanceLogs
	if err := query.
		Preload("User").
		Order("clock_in_time DESC, id DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal mengambil riwayat absensi"})
	}

	var result []responses.LogsRes
	for _, l := range logs {
		result = append(result, mappers.ToLogsResponse(l))
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"logs":       result,
			"totalPages": totalPages,
			"page":       page,
			"limit":      limit,
			"total":      total,
		},
	})
}
