package query

import (
	"github.com/CallumLewisGH/Generic-Service-Base/database"
)

func CheckDatabaseHealth() string {
	if database.GetDatabase().CheckDatabaseHealth() != nil {
		return "database unhealthy! x_x!"
	}
	return "database healthy! :)"
}
