package services

import (
	"fmt"
	"github.com/xuri/excelize/v2"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/repo"
)

func GenerateAttendanceExcel(kelas, jurusan, tanggal string) (*excelize.File, error) {

	rows, err := repo.GetAttendanceRows(kelas, jurusan, tanggal)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()

	sheet := "Absensi"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// =====================
	// STYLE
	// =====================
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})

	hadirStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#C6EFCE"}, Pattern: 1},
	})

	telatStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFD966"}, Pattern: 1},
	})

	absenStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#F8CBAD"}, Pattern: 1},
	})

	// =====================
	// HEADER
	// =====================
	headers := []string{"No", "NISN", "Nama", "Kelas", "Status", "Waktu"}

	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// =====================
	// DATA + SUMMARY
	// =====================
	hadir := 0
	telat := 0
	belum := 0

	for i, r := range rows {

		row := i + 2

		status := "-"
		waktu := "-"

		style := absenStyle

		if r.Status != nil && r.ClockInTime != nil {
			status = *r.Status
			waktu = r.ClockInTime.Format("2006-01-02 15:04:05")

			if status == "hadir" {
				hadir++
				style = hadirStyle
			} else {
				telat++
				style = telatStyle
			}
		} else {
			belum++
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.Nisn)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.FullName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), r.ClassGroup)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), status)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), waktu)

		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("F%d", row), style)
	}

	// =====================
	// SHEET SUMMARY
	// =====================
	summary := "Summary"
	f.NewSheet(summary)

	f.SetCellValue(summary, "A1", "Total Hadir")
	f.SetCellValue(summary, "B1", hadir)

	f.SetCellValue(summary, "A2", "Total Telat")
	f.SetCellValue(summary, "B2", telat)

	f.SetCellValue(summary, "A3", "Belum Absen")
	f.SetCellValue(summary, "B3", belum)

	f.SetCellValue(summary, "A4", "Total Siswa")
	f.SetCellValue(summary, "B4", len(rows))

	return f, nil
}
