package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"player_rest_server/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	server.GET("/api/players/:player_id", handlePlayer)
	server.GET("/api/players", handlePlayers)
}

func handlePlayers(ctx *gin.Context) {
	var err error
	page := -1
	pageSize := 250

	fmt.Println("INFO: Handling request for all players")
	pageStr := ctx.Query("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{"error": fmt.Sprintf("invalid value for page: %s", pageStr)})
		}
	}

	if page < 0 {
		ctx.JSON(http.StatusExpectationFailed, gin.H{"error": fmt.Sprintf("invalid value for page: %d", page)})
	}

	pageSizeStr := ctx.Query("page_size")
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			ctx.JSON(http.StatusExpectationFailed, gin.H{"error": fmt.Sprintf("invalid value for page_size: %s", pageSizeStr)})
		}
	}

	if pageSize < 0 || pageSize > 10000 {
		ctx.JSON(http.StatusExpectationFailed, gin.H{"error": fmt.Sprintf("invalid value for page_size: %d", pageSize)})
	}

	fmt.Println("INFO: page:", page, "pageSize:", pageSize)

	players, err := db.FetchPlayers(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, players)
}

func handlePlayer(ctx *gin.Context) {
	playerId := ctx.Param("player_id")
	fmt.Println("INFO: Handling request for player ID", playerId)
	player, err := db.FetchPlayer(playerId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if player == nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, player)
}
