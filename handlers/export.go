package handlers

import (
	"fmt"
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func ExportAttendance(c *fiber.Ctx) error {
	kelas := c.Query("kelas")
	jurusan := c.Query("jurusan")
	tanggal := c.Query("tanggal")
	var logs []models.AttedanceLogs
	query := database.DB.Preload("User").Preload("Token")

	if kelas != "" {
		query = query.Joins("JOIN users ON users.id = attedance_logs.user_id").
			Where("users.class_group = ?", kelas)
	}

	if jurusan != "" {
		if kelas == "" {
			query = query.Joins("JOIN users ON users.id = attedance_logs.user_id")
		}
		query = query.Where("users.class_group LIKE ?", "%"+jurusan+"%")
	}

	if tanggal != "" {
		query = query.Where("DATE(clock_in_time) = ?", tanggal)
	}

	if err := query.Find(&logs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal mengambil data absensi"})
	}

	f := excelize.NewFile()
	sheet := "Absensi"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#DDDDDD"}, Pattern: 1},
	})

	headers := []string{"No", "NISN", "Nama Lengkap", "Kelas", "Status", "Waktu Absen"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 5)
	f.SetColWidth(sheet, "B", "B", 15)
	f.SetColWidth(sheet, "C", "C", 25)
	f.SetColWidth(sheet, "D", "D", 15)
	f.SetColWidth(sheet, "E", "E", 10)
	f.SetColWidth(sheet, "F", "F", 20)

	for i, log := range logs {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), log.User.Nisn)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), log.User.FullName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), log.User.ClassGroup)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), log.Status)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), log.ClockInTime.Format("2006-01-02 15:04:05"))
	}

	fileName := fmt.Sprintf("absensi_%s.xlsx", time.Now().Format("2006-01-02"))

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	return f.Write(c.Response().BodyWriter())
}