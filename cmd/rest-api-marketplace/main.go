package main

import (
	_ "github.com/go-playground/validator/v10"
	_ "github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"rest-api-marketplace/internal/app"
)

func main() {
	app.Run()
}
