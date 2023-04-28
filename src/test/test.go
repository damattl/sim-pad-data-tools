package main

import (
	"fmt"
	"time"
)

func main() {
	var dateTime = "2023-04-13T07:38:14"
	parsedTime, _ := time.Parse("2006-01-02T15:04:05", dateTime)
	fmt.Println(parsedTime.Format("2006-01-02"))
}
