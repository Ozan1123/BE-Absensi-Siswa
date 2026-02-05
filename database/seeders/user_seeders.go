package seeders

import (
	"fmt"
	"log"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"

	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

func SeedUsersFromExcel(path string) error {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return err
	}

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(row[3]), bcrypt.DefaultCost)

		user := models.Users{
			Nisn:       row[0],
			FullName:   row[1],
			Username:   row[2],
			Password:   string(passwordHash),
			Role:       "siswa",
			ClassGroup: row[4],
		}

		// hindari duplicate username
		var existing models.Users
		err := database.DB.Where("username = ?", user.Username).First(&existing).Error

		if err == nil {
			log.Println("skip duplicate:", user.Username)
			continue
		}

		if err := database.DB.Create(&user).Error; err != nil {
			log.Println("insert failed:", err)
			continue
		}

		fmt.Println("inserted:", user.Username)
	}

	fmt.Println("Users seeded from Excel selesai")
	return nil
}
