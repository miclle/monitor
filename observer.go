package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/miclle/observer/config"
	"github.com/miclle/observer/detector"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)
}

func main() {

	detector.Init(config.Mongo.Host, config.Mongo.Name, config.Mongo.Mode)

	router := gin.Default()
	// Monitor
	router.GET("/tasks", func(c *gin.Context) {
		tasks, err := detector.TaskMgr.List()

		if err != nil {
			log.Error(err.Error())
			c.JSON(400, err)
			c.Abort()
			return
		}
		c.JSON(200, gin.H{"tasks": tasks})
	})

	router.Run(fmt.Sprintf(":%d", config.Port))
}
