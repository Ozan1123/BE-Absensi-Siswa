package services

import (
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/repo"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/utils"
)

func GenerateAttendanceExcel(kelas, jurusan, startDate, endDate string) (*excelize.File, error) {

	users, err := repo.GetAttendanceRows(kelas, jurusan, startDate, endDate)
	if err != nil {
		return nil, err
	}

	start, end, err := utils.DateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Generate all dates in range
	var dateList []time.Time
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateList = append(dateList, d)
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
	sakitStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#BDD7EE"}, Pattern: 1},
	})

	// =====================
	// HEADER
	// =====================
	colIndex := 1
	rowIndex := 1

	staticHeaders := []string{"No", "NISN", "Nama", "Kelas"}
	for _, h := range staticHeaders {
		cell, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		colIndex++
	}

	for _, d := range dateList {
		cell, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cell, d.Format("02-Jan"))
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		colIndex++
	}

	summaryHeaders := []string{"Total Hadir", "Total Sakit", "Total Izin", "Total Telat", "Total Alfa"}
	for _, h := range summaryHeaders {
		cell, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		colIndex++
	}

	// =====================
	// DATA + SUMMARY
	// =====================
	globalHadir, globalTelat, globalAlfa, globalSakit, globalIzin, globalBelum := 0, 0, 0, 0, 0, 0

	rowIndex = 2
	for i, u := range users {
		colIndex = 1

		// Static Info
		cellNo, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cellNo, i+1)
		colIndex++

		cellNisn, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cellNisn, u.Nisn)
		colIndex++

		cellNama, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cellNama, u.FullName)
		colIndex++

		cellKelas, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cellKelas, u.ClassGroup)
		colIndex++

		// Map Logs
		logMap := make(map[string]string)
		for _, log := range u.AttedanceLogs {
			if !log.ClockInTime.IsZero() {
				dateStr := log.ClockInTime.Format("2006-01-02")
				logMap[dateStr] = log.Status
			}
		}

		userHadir, userSakit, userIzin, userTelat, userAlfa := 0, 0, 0, 0, 0

		for _, d := range dateList {
			dateStr := d.Format("2006-01-02")
			status, exists := logMap[dateStr]
			if !exists {
				status = "belum_absen"
			}

			style := absenStyle
			switch status {
			case "hadir":
				userHadir++
				globalHadir++
				style = hadirStyle
			case "telat":
				userTelat++
				globalTelat++
				style = telatStyle
			case "alfa":
				userAlfa++
				globalAlfa++
				style = absenStyle
			case "sakit":
				userSakit++
				globalSakit++
				style = sakitStyle
			case "izin":
				userIzin++
				globalIzin++
				style = sakitStyle
			default:
				globalBelum++
			}

			cell, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
			f.SetCellValue(sheet, cell, status)
			f.SetCellStyle(sheet, cell, cell, style)
			colIndex++
		}

		// Write user totals
		cellHadir, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cellHadir, userHadir)
		colIndex++

		cellSakit, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cellSakit, userSakit)
		colIndex++

		cellIzin, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cellIzin, userIzin)
		colIndex++

		cellTelat, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cellTelat, userTelat)
		colIndex++

		cellAlfa, _ := excelize.CoordinatesToCellName(colIndex, rowIndex)
		f.SetCellValue(sheet, cellAlfa, userAlfa)
		colIndex++

		rowIndex++
	}

	// =====================
	// SHEET SUMMARY
	// =====================
	summary := "Summary"
	f.NewSheet(summary)

	f.SetCellValue(summary, "A1", "Total Hadir")
	f.SetCellValue(summary, "B1", globalHadir)

	f.SetCellValue(summary, "A2", "Total Telat")
	f.SetCellValue(summary, "B2", globalTelat)

	f.SetCellValue(summary, "A3", "Total Alfa")
	f.SetCellValue(summary, "B3", globalAlfa)

	f.SetCellValue(summary, "A4", "Total Sakit")
	f.SetCellValue(summary, "B4", globalSakit)

	f.SetCellValue(summary, "A5", "Total Izin")
	f.SetCellValue(summary, "B5", globalIzin)

	f.SetCellValue(summary, "A6", "Belum Absen")
	f.SetCellValue(summary, "B6", globalBelum)

	f.SetCellValue(summary, "A7", "Total Siswa")
	f.SetCellValue(summary, "B7", len(users))

	return f, nil
}

