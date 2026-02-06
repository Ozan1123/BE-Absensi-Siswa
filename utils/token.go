package utils

import (
	"errors"
	"math/rand"
	"time"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandomString(lenght int) string {
	b := make([]byte, lenght)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func CreateToken(adminID int64, durationMinutes int, lateAfter int) (*models.AttedanceTokens, error) {
	var tokenCode string

	for {
		code := RandomString(6)

		var count int64
		database.DB.Model(&models.AttedanceTokens{}).
			Where("token_code = ?", code).
			Count(&count)
		
		if count == 0 {
			tokenCode = code
			break
		}
	}
	
	
	token := models.AttedanceTokens{
		TokenCode:  tokenCode,
		CreatedBy:  adminID,
		IsActive:   true,
		LateAfter: time.Now().Add(time.Minute * time.Duration(lateAfter)),
		ValidUntil: time.Now().Add(time.Minute * time.Duration(durationMinutes),),
	}

	if err := database.DB.Create(&token).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("User").First(&token, token.ID).Error; err != nil {
		return nil, err
	}

	return &token, nil
}


func VerifyTokenCode(input string) (*models.AttedanceTokens, error) {
	var token models.AttedanceTokens

	err := database.DB.
		Where("token_code = ? AND is_active = ?", input, true).
		First(&token).Error
	if err != nil {
		return  nil, errors.New("Token tidak ditemukan!")
	}

	// cek expired
	if time.Now().After(token.ValidUntil) {
		err := database.DB.
			Model(&models.AttedanceTokens{}).
			Where("id = ?", token.ID).
			Update("is_active", false).Error
		if err != nil {
			return  nil, err
		}

		return  nil, errors.New("Token sudah expired!")
	}

	return &token, nil
}

