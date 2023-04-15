package main

import "fmt"

func main() {
	fmt.Println("Starting with parsing")
	err := parse()
	if err != nil {
		fmt.Println("Following Error occurred: ")
		fmt.Println(err)
		_, _ = fmt.Scanln()
		fmt.Println("Press any key to quit")
		return
	}
	fmt.Println("File generated")

	fmt.Println("Press any key to quit")
	_, _ = fmt.Scanln()
}
