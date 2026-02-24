package services

import (
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
