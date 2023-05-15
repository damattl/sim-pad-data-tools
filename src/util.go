package main

import (
	"fmt"
	"log"
	"os"
)

func printErrorAndWaitForExit(err error) {
	log.Println("Following Error occurred: ")

	log.Println(err)
	_, _ = fmt.Scanln()
	fmt.Println("Press enter to quit")
	log.Println("Program closed with error")
	os.Exit(1)
}
