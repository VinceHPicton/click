package db

// const (
// 	userColumns = `(
// 		name,
// 		bio,
// 		birth_date,
// 		last_location_long,
// 		mobile,
// 		last_active,
// 		email,
// 		sex,
// 		interested_in
// 	)`
// )

// func (dba *DBAccess) CreateUser() error {

// 	query := fmt.Sprintf("INSERT INTO users %v VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", userColumns)

// 	statement, err := dba.DB.Prepare(query)
// 	if err != nil {
// 		return err
// 	}
// 	defer statement.Close()

// 	_, err = statement.Exec(1, 1, 1, "2022-10-10 11:30:30")
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
