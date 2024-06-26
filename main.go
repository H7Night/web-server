package main

import (
	"web-server/models"
	"web-server/routes"
)

func main() {
	models.InitDB()
	routes.InitRouter()
}
