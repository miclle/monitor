package main

import (
	"github.com/gin-gonic/gin"
	"github.com/miclle/observer/detector"
	"qiniupkg.com/x/log.v7"
)

func main() {
	router := gin.Default()

	host := "127.0.0.1:27017"
	name := "observer_test"
	mode := "strong"

	detector.Init(host, name, mode)

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

	router.Run(":8000")
}
