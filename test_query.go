//go:build ignore

package main

import (
	"fmt"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/config"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/repo"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()
	rows, err := repo.GetAttendanceRows("", "", "2024-07-01", "2024-07-31")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Total rows:", len(rows))
	
	counts := make(map[string]int)
	for _, r := range rows {
		counts[r.FullName]++
	}
	
	hasMultiple := false
	for name, count := range counts {
		if count > 1 {
			fmt.Printf("Student %s has %d rows\n", name, count)
			hasMultiple = true
			break
		}
	}
	if !hasMultiple {
		fmt.Println("NO student has multiple rows!")
	}
}
