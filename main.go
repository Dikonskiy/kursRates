package main

import (
	"kursRates/internal/app"

	_ "kursRates/docs"

	_ "github.com/go-sql-driver/mysql"
)

// @title Swagger kursRates API
// @version 0.1
// @description A web service that, upon request, collects data from the public API of the national bank and saves the data to the local TEST database
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	app := app.NewApplication()

	//commit 1
	//commit 2
	app.StartServer()
}
