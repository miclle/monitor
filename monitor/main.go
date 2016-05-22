package main

import (
	"fmt"

	"github.com/carlescere/scheduler"
)

func main() {
	job := func() {
		fmt.Println("Time's up!")
	}
	scheduler.Every(5).Seconds().Run(job)
	scheduler.Every().Day().Run(job)
	scheduler.Every().Sunday().At("08:30").Run(job)
}
