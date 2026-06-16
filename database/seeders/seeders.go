package seeders

func RunSeed() error {
	if err := SeedUsersFromExcel("database/seeders/files/admin-Superadmin.xlsx"); err != nil {
		return err
	}

	// Seed default notification settings (WA config) jika tabel kosong
	if err := SeedNotificationSettings(); err != nil {
		return err
	}

	return nil
}





