package main

import "fmt"

func main() {
	fmt.Println("Starting with parsing")
	parse()
	fmt.Println("File generated")

	fmt.Println("Press any key to quit")
	_, _ = fmt.Scanln()
}
