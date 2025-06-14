package main

import (
	"fmt"
	"net/http"

	"github.com/CallumLewisGH/Generic-Service-Base/database"
	_ "github.com/CallumLewisGH/Generic-Service-Base/docs"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api/authentication"
	"github.com/CallumLewisGH/Generic-Service-Base/internal/api/routes"
	_ "github.com/joho/godotenv"
	_ "github.com/lann/builder"
)

func main() {
	//Get Authentication Config
	authentication.SetupGoogleOAuth()

	//Database Initialisation
	database.GetDatabase()

	//Creates new server instance
	srv := api.NewServer()

	//Route Registry
	routes.RegisterDatabaseRoutes(srv)
	routes.RegisterUserRoutes(srv)
	routes.RegisterAuthenticationRoutes(srv)

	println(fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", "http://localhost:8080/swagger/index.html", "SwaggerUI: http://localhost:8080/swagger/index.html"))

	//Starts the server on port 8080
	http.ListenAndServe(":8080", srv)

}

//Implementation Notes
// Queries IE repo.find, .first, all accept pointers &user and data is returned to that object with an err being returned from the function
// Commands IE repo.CreateOne, UpdateMany, DeleteOne, all accept actual values this way the input value is discaurded and another value is supplimented

//COMMANDS THAT ARE GOOD TO HAVE EASY ACCESS TO
// - docker compose up
// - swag init -g cmd/app/main.go --dir ./
// - go run cmd/app/main.go
// - go test ./unit_tests

// TODO:
// Role Based Access control IE Authroisation => Where to apply authroisation middleware? on the commands?
//  - roles permissions ect
//    - Should be able to create permissions for each command and then apply based off of that
//    - There should be a way to scope your database write access to fields only associated with you

// Make commands and queries asynchronous using go routines:
//  - This will speed things up significantly

// IntegrationTests should work -> Spin up docker containers per each one

// Look into configuring air for hot reloads
