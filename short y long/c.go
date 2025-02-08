package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var dataVersion int = 0
var mu sync.Mutex

func main() {
	router := gin.Default()

	// Endpoint para recibir cambios
	router.POST("/update", func(ctx *gin.Context) {
		mu.Lock()
		dataVersion++
		mu.Unlock()

		ctx.JSON(http.StatusOK, gin.H{"message": "Data updated"})
	})

	// Endpoint Long Polling
	router.GET("/wait_changes/:version", func(ctx *gin.Context) {
		clientVersion := ctx.Param("version")

		// Espera cambios antes de responder
		for {
			mu.Lock()
			currentVersion := dataVersion
			mu.Unlock()

			if fmt.Sprintf("%d", currentVersion) != clientVersion {
				ctx.JSON(http.StatusOK, gin.H{"new_version": currentVersion})
				return
			}

			time.Sleep(500 * time.Millisecond) // Pequeña pausa antes de volver a revisar
		}
	})

	router.Run(":5000")
}