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

var (
	hello string = "hello"
)

func main() {
	//Get Authentication Config
	authentication.SetupGoogleOAuth()

	//Database Initialisation
	database.SetDevMode() // Change to prod mode during production
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
// New models need to be added to the model_registry in the database package
// All usefull commands have now been moved to the justfile
// Change to prod mode for the database when deploying to prod => Also change GIN server from default

// TODO:

// Add rate limiting for security

// Look at crfs => something like that security issues

// Find a way to make the API private as to not accept requests from anywhere but a specified frontend domain
//  - This is super important for me to maintain security I don't think I can code well enough right now to have a public facing api
//  -

// Role Based Access control IE Authroisation => Where to apply authroisation middleware? on the commands?
//  - roles permissions ect
//    - Should be able to create permissions for each command and then apply based off of that
//    - There should be a way to scope your database write access to fields only associated with you
