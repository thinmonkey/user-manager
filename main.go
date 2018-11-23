package main

import (
	"github.com/thinmonkey/user-manager/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/thinmonkey/user-manager/router"
	"github.com/thinmonkey/user-manager/utils/log"
)

func main() {

	config.Init()

	// Create the Gin engine.
	gin := gin.New()

	// Routes.
	router.Load(
		// Cores.
		gin,
		// Middlwares.
	)

	log.Info(http.ListenAndServe(":8088", gin).Error())
}
