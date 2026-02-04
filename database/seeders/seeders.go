package seeders

func RunSeed() error {
	if err := SeedUsersFromExcel("database/seeders/files/testing.xlsx"); err != nil {
		return err
	}

	return nil
}





