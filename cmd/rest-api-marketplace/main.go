// General package of service
// main.go is actually the main entry point of web service
package main

import (
	"rest-api-marketplace/internal/app"

	_ "github.com/go-playground/validator/v10"
	_ "github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	app.Run()
}
