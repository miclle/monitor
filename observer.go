package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/renders/multitemplate"
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
	render := multitemplate.New()

	templateBox := rice.MustFindBox("templates")

	templateBox.Walk("", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			templateString := templateBox.MustString(path)
			tmplMessage, err := template.New(path).Parse(templateString)
			if err != nil {
				return err
			}
			render.Add(path, tmplMessage)
		}
		return nil
	})

	router.HTMLRender = render

	router.StaticFS("/assets", rice.MustFindBox("assets").HTTPBox())

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

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
