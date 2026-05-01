package requests

// UpdateSettingReq untuk update single setting
type UpdateSettingReq struct {
	SettingKey   string `json:"setting_key" validate:"required"`
	SettingValue string `json:"setting_value" validate:"required"`
}

// UpdateSettingsBulkReq untuk update multiple settings sekaligus
type UpdateSettingsBulkReq struct {
	Settings []UpdateSettingReq `json:"settings" validate:"required"`
}

// TestWAReq untuk test kirim WA
type TestWAReq struct {
	Phone   string `json:"phone" validate:"required"`
	Message string `json:"message" validate:"required"`
}

// UpdateStudentStatusReq untuk guru/admin set status siswa
type UpdateStudentStatusReq struct {
	UserID int64  `json:"user_id" validate:"required"`
	Status string `json:"status" validate:"required"` // sakit, izin, alfa
}
