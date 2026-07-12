package responses

import "time"

type LogsRes struct {
	ID          int64     `json:"id"`
	User        UserMini  `json:"user"`
	Status      string    `json:"status"`
	CapturedIp  *string   `json:"captured_ip"`
	ClockInTime time.Time `json:"clock_in_time"`
}

type LogResMini struct {
	ID          uint       `json:"id"`
	Status      string     `json:"status"`
	CapturedIP  *string    `json:"captured_ip"`
	ClockInTime time.Time  `json:"clock_in_time"`
}

type LogResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
