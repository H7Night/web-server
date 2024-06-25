package main

import (
	"goIland/models"
	"goIland/routes"
)

func main() {
	models.InitDB()
	routes.InitRouter()
}
