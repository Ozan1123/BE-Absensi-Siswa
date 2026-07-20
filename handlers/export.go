package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/services"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
)

// ExportAttendance godoc
// @Summary Export data absensi ke Excel
// @Description Mengexport data absensi berdasarkan kelas, jurusan, dan rentang tanggal dalam bentuk file Excel
// @Tags attendance
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param kelas query string false "Kelas"
// @Param jurusan query string false "Jurusan"
// @Param start_date query string true "Tanggal mulai (format: YYYY-MM-DD)"
// @Param end_date query string true "Tanggal akhir (format: YYYY-MM-DD)"
// @Success 200 {file} file
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /export/attendance [get]
func ExportAttendance(c *fiber.Ctx) error {

	kelas := c.Query("kelas")
	jurusan := c.Query("jurusan")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	start, end, err := utils.DateRange(startDate, endDate)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	file, err := services.GenerateAttendanceExcel(kelas, jurusan, startDate, endDate)
	if err != nil {
		return err
	}

	var fileName string
	if startDate == endDate || endDate == "" {
		fileName = fmt.Sprintf("absensi_%s.xlsx", start.Format("2006-01-02"))
	} else {
		fileName = fmt.Sprintf("absensi_%s_sampai_%s.xlsx", start.Format("2006-01-02"), end.Format("2006-01-02"))
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+fileName)

	return file.Write(c.Response().BodyWriter())
}
