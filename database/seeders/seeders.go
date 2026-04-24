package seeders

func RunSeed() error {
	if err := SeedUsersFromExcel("database/seeders/files/admin-Superadmin.xlsx"); err != nil {
		return err
	}
	return nil
}





