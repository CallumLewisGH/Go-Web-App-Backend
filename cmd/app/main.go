package main

import (
	"fmt"
	"net/http"

	"github.com/CallumLewisGH/Generic-Service-Base/database"
	_ "github.com/CallumLewisGH/Generic-Service-Base/docs"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api/routes"
	_ "github.com/joho/godotenv"
	_ "github.com/lann/builder"
)

func main() {
	//Database Initialisation
	database.GetDatabase()

	//Creates new server instance
	srv := api.NewServer()

	//Route Registry
	routes.RegisterDatabaseRoutes(srv)
	routes.RegisterUserRoutes(srv)

	println(fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", "http://localhost:8080/swagger/index.html", "SwaggerUI: http://localhost:8080/swagger/index.html"))

	//Starts the server on port 8080
	http.ListenAndServe(":8080", srv)

}

//TODO for full architecture
// BARE BONES:
//  - HTTP Routing CHECK
//  - Database connection CHECK
//  - Object mapping CHECK
//  - Database interactions CHECK
//  - Service, data and presentation layers CHECK

// Do just to see what challenges you aren't aware of:
//  - Create Delete UserById route CHECK
//  - Create Update UserById route CHECK
//  - Create Get    UserById route CHECK

// WANT TO HAVE:
//  - QueryStringBuilder type thing CHECK (THAT NEED HAS BEEN FILLED)
//  - Authentication HALFWAY => I think I can use GoAuth with a auth provider IE google
//     - I can implement Authentication Middleware on routes
//     - I can foresee a problem with RBAC where I can't call more than one command/query from the same route. Since I won't be able to effectively "rollback" the transactions
//  - Forgot my passcode
//  - 2FA
//  - Role based access control NOT EVEN CLOSE

//Implementation Notes
// Queries IE repo.find, .first, all accept pointers &user and data is returned to that object with an err being returned from the function
// Commands IE repo.CreateOne, UpdateMany, DeleteOne, all accept actual values this way the input value is discaurded and another value is supplimented

//COMMANDS THAT ARE GOOD TO HAVE EASY ACCESS TO
// - docker compose up
// - swag init -g cmd/app/main.go --dir ./
// - go run cmd/app/main.go
//
