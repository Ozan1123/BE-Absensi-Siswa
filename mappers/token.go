package mappers

import (
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/responses"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"
)

func ToTokenResponse(t *models.AttedanceTokens) responses.TokenRes {
	return responses.TokenRes{
		ID: t.ID,
		TokenCode: t.TokenCode,
		CreatedBy: responses.UserMini{
			ID: t.User.ID,
			FullName: t.User.FullName,
		},
		IsActive: t.IsActive,
		ValidUntil: t.ValidUntil,
		CreatedAt: t.CreatedAt,
	}
}