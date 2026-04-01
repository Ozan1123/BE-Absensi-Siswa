package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/services"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
)


// ExportAttendance godoc
// @Summary Export data absensi ke Excel
// @Description Mengexport data absensi berdasarkan kelas, jurusan, dan tanggal dalam bentuk file Excel
// @Tags attendance
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param kelas query string false "Kelas"
// @Param jurusan query string false "Jurusan"
// @Param tanggal query string true "Tanggal (format: YYYY-MM-DD)"
// @Success 200 {file} file
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /export/attendance [get]
func ExportAttendance(c *fiber.Ctx) error {

	kelas := c.Query("kelas")
	jurusan := c.Query("jurusan")
	tanggal := c.Query("tanggal")

	start, _, err := utils.DayRange(tanggal)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "format tanggal salah (YYYY-MM-DD)"})
	}

	file, err := services.GenerateAttendanceExcel(kelas, jurusan, tanggal)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("absensi_%s.xlsx", start.Format("2006-01-02"))

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+fileName)

	return file.Write(c.Response().BodyWriter())
}

