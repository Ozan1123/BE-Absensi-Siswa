================================================================================================================= To Install Package Migration : go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

================================================================================================================= To Create Migration : migrate create -ext sql -dir database/migration -seq create_examples_table

================================================================================================================= To Running (UP) Migration (Example Mysql) : migrate -database "mysql://db_user:db_password@tcp(127.0.0.1:3306)/db_name" -path database/migration up

================================================================================================================== To Running (DOWN) Migration : migrate -database "mysql://db_user:db_password@tcp(127.0.0.1:3306)/db_name" -path database/migration down