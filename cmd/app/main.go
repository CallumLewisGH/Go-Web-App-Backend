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

// WANT TO HAVE:
//  - QueryStringBuilder
//  - Authentication
//  - Role based access control
