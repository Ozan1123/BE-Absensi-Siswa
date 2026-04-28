package services

import (
	"fmt"

	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/dto/responses"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/models"

	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)



func ImportUsersFromExcel(path string) (*responses.ImportResult,error){

	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	result := &responses.ImportResult{}

	for i,row := range rows {

		if i == 0 {
			continue
		}

		if len(row) < 5 {
			result.Failed++
			continue
		}

		passwordHash,_ := bcrypt.GenerateFromPassword(
			[]byte(row[3]),
			bcrypt.DefaultCost,
		)

		user := models.Users{
			Nisn: row[0],
			FullName: row[1],
			Username: row[2],
			Password: string(passwordHash),
			Role: row[5],
			ClassGroup: row[4],
			ParentPhone: row[6],
		}

		var existing models.Users

		err := database.DB.
			Where("username = ?",user.Username).
			First(&existing).Error

		if err == nil {
			result.Duplicates++
			result.SkippedUsers = append(
				result.SkippedUsers,
				user.Username,
			)
			continue
		}

		if err := database.DB.Create(&user).Error; err != nil{
			result.Failed++
			continue
		}

		result.Inserted++
		fmt.Println("Inserted:",user.Username)
	}

	return result,nil
}