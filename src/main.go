package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func logSetup(logDir string) {
	err := os.MkdirAll(logDir, 0777)
	if err != nil {
		fmt.Printf("could not create log dir due to err: %v \n", err)
		_, _ = fmt.Scanln()
		fmt.Println("Press enter to quit")
		os.Exit(1)
	}
	now := time.Now()
	f, err := os.OpenFile(logDir+"/log-"+now.Format("02-01-2006-15-04-05"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		fmt.Printf("could not open a log file due to err: %v \n", err)
		_, _ = fmt.Scanln()
		fmt.Println("Press enter to quit")
		os.Exit(1)
	}

	defer f.Close()

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
}

func inputSetup(inputDir string) {
	err := os.MkdirAll(inputDir, 0777)
	if err != nil {
		fmt.Printf("could not create input dir due to err: %v \n", err)
		_, _ = fmt.Scanln()
		fmt.Println("Press enter to quit")
		os.Exit(1)
	}
	fmt.Println(os.Executable())
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("could not read input files due to err: %v \n", err)
		_, _ = fmt.Scanln()
		fmt.Println("Press enter to quit")
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("There are currently no files in the input dir, please put them there.")

		fmt.Println("To continue press enter.")
		fmt.Println("To exit, enter 'exit' and press enter")
		var scan string
		_, _ = fmt.Scanln(&scan)
		if strings.ToUpper(scan) == "EXIT" {
			os.Exit(0)
		}
	}
}

func main() {
	execPath, err := os.Executable()
	if err != nil {
		printErrorAndWaitForExit(err)
	}
	execDir := path.Dir(execPath)
	if err != nil {
		printErrorAndWaitForExit(err)
	}

	inputDir := path.Join(execDir, "input")

	logSetup(path.Join(execDir, "logs"))
	inputSetup(inputDir)

	log.Println("Starting with parsing")
	err = parse(inputDir, execDir)
	if err != nil {
		printErrorAndWaitForExit(err)
	}
	log.Println("File generated")

	fmt.Println("Press enter to quit")
	_, _ = fmt.Scanln()
	log.Println("Program closed successfully")
}
