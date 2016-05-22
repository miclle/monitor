package main

import (
	"fmt"

	"github.com/carlescere/scheduler"
)

func main() {
	job := func() {
		fmt.Println("Time's up!")
	}
	scheduler.Every(5).Seconds().Run(function)
	scheduler.Every().Day().Run(function)
	scheduler.Every().Sunday().At("08:30").Run(function)
}
