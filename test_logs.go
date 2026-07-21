//go:build ignore

package main

import (
	"fmt"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/config"
	"github.com/KicauOrgspark/BE-Absensi-Siswa/database"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()
	var count int64
	database.DB.Table("attedance_logs").Count(&count)
	fmt.Println("Total logs in database:", count)
	
	type Log struct {
	    UserID int
	    ClockInTime string
	}
	var logs []Log
	database.DB.Table("attedance_logs").Select("user_id, clock_in_time").Scan(&logs)
	for _, l := range logs {
	    fmt.Printf("User %d: %s\n", l.UserID, l.ClockInTime)
	}
}
