package main

import (
	"fmt"
	"time"

	"github.com/carlescere/scheduler"
)

func main() {
	fmt.Println("Main execute at:\t", time.Now())

	scheduler.Every(5).Seconds().Run(func() {
		fmt.Println("Execute every 5 seconds\t", time.Now())
	})

	scheduler.Every(1).Minutes().Run(func() {
		fmt.Println("Execute every 1 minute\t", time.Now())
	})

	scheduler.Every().Day().Run(func() {
		fmt.Println("Execute every 1 day\t", time.Now())
	})

	scheduler.Every().Sunday().At("08:30").Run(func() {
		fmt.Println("Execute every Sunday 08:30\t", time.Now())
	})

	done := make(chan struct{})

	<-done
}
