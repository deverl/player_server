package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"player_server/db"
	"player_server/routes"

	"github.com/gin-gonic/gin"
)

const (
	csvRelativePath = "./csv/Player.csv"
	port            = 8800
)

func main() {
	fmt.Println("INFO: Starting api server")

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("ERROR: Couldn't get working directory. err:", err)
	}

	fmt.Printf("INFO: Working directory: '%s'\n", dir)

	csvPath := filepath.Join(dir, csvRelativePath)

	fmt.Printf("INFO: csvPath: '%s'\n", csvPath)

	if db.GetDB() == nil {
		fmt.Println("ERROR: Database connection not ready... exiting")
		return
	}

	fmt.Println("INFO: Populating the player database")

	db.PopulatePlayer(csvPath)

	server := gin.Default()

	routes.RegisterRoutes(server)

	go populateDB(csvPath)

	serverString := fmt.Sprintf(":%d", port)

	server.Run(serverString)
}

func populateDB(csvPath string) {
	// Check every minute.
	time.Sleep(60 * time.Second)
	db.PopulatePlayer(csvPath)
}
