package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	err := os.MkdirAll("logs", 0777)
	if err != nil {
		fmt.Printf("could not create log dir due to err: %v \n", err)
		_, _ = fmt.Scanln()
		fmt.Println("Press any key to quit")
		os.Exit(1)
	}
	now := time.Now()
	f, err := os.OpenFile("logs/log-"+now.Format("02-01-2006-15-04-05"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		fmt.Printf("could not open a log file due to err: %v \n", err)
		_, _ = fmt.Scanln()
		fmt.Println("Press any key to quit")
		os.Exit(1)
	}

	defer f.Close()

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)

	log.Println("Starting with parsing")
	err = parse()
	if err != nil {
		log.Println("Following Error occurred: ")

		log.Println(err)
		_, _ = fmt.Scanln()
		fmt.Println("Press any key to quit")
		log.Println("Program closed with error")
		return
	}
	log.Println("File generated")

	fmt.Println("Press any key to quit")
	_, _ = fmt.Scanln()
	log.Println("Program closed successfully")
}
