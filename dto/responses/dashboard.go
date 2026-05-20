package responses

type DashboardResponse struct {
	TotalTokens       int `json:"total_token"`
	TokenHariIni      int `json:"token_hari_ini"`
	ActiveTokens      int `json:"token_aktif"`
	TotalAbsenHariIni int `json:"total_absen_hari_ini"`
	TotalHadir        int `json:"total_hadir_hari_ini"`
	TotalTelat        int `json:"total_telat_hari_ini"`
}