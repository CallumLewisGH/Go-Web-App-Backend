package command

import "github.com/CallumLewisGH/Generic-Service-Base/database"

func RunDatabaseMigrations() error {
	return database.GetDatabase().RunMigrations()
}
