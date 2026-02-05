package utils

import (
	"math/rand"
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
)

const charset = "QMBLO01386"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func CreateToken(adminID int64, durationMinutes int) (*models.AttedanceTokens, error) {
	token := models.AttedanceTokens{
		TokenCode:  RandomString(10),
		CreatedBy:  adminID,
		IsActive:   true,
		ValidUntil: time.Now().Add(time.Minute * time.Duration(durationMinutes)),
	}

	if err := database.DB.Create(&token).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("User").First(&token, token.ID).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func VerifyTokenCode(input string) (bool, error) {
	var token models.AttedanceTokens

	// ambil token dari DB
	err := database.DB.
		Where("token_code = ? AND is_active = ?", input, true).
		First(&token).Error
	if err != nil {
		return false, err
	}

	// cek expired
	if time.Now().After(token.ValidUntil) {
		// nonaktifkan token
		err := database.DB.
			Model(&models.AttedanceTokens{}).
			Where("id = ?", token.ID).
			Update("is_active", false).Error
		if err != nil {
			return false, err
		}

		return false, nil // ❗ expired = false
	}

	return true, nil
}

