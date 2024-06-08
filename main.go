package main

import (
	"fmt"
	"time"

	"player_server/db"
	"player_server/routes"

	"github.com/gin-gonic/gin"
)

const (
	csvPath = "./csv/Player.csv"
	port    = 8800
)

func main() {
	fmt.Println("INFO: Starting api server")

	if db.GetDB() == nil {
		fmt.Println("ERROR: Database connection not ready... exiting")
		return
	}

	fmt.Println("INFO: Populating the player database")

	db.PopulatePlayer(csvPath)

	server := gin.Default()

	routes.RegisterRoutes(server)

	go populateDB()

	serverString := fmt.Sprintf(":%d", port)

	server.Run(serverString)
}

func populateDB() {
	time.Sleep(3600 * time.Hour)
	db.PopulatePlayer(csvPath)
}
