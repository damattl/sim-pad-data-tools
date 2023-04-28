package main

import (
	"encoding/xml"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"sort"
	"time"
)

const logPath = "/EventLog.xml"
const cprEventsPath = "/CPR/CPREvents.xml"

func readAndParseFiles(inputDir string) ([]SimPadData, error) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return []SimPadData{}, err
	}

	var results []SimPadData

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		rawCPRData, err := os.ReadFile(inputDir + "/" + entry.Name() + cprEventsPath)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		rawLogData, err := os.ReadFile(inputDir + "/" + entry.Name() + logPath)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		var parsedLog SimPadLog
		err = xml.Unmarshal(rawLogData, &parsedLog)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		var parsedList SimPadCPREventList
		err = xml.Unmarshal(rawCPRData, &parsedList)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		results = append(results, SimPadData{
			EventList: parsedList,
			Log:       parsedLog,
		})

	}
	return results, err
}

func createCPRParamMap() map[string]SimPadCPREventParameter {
	paramMap := map[string]SimPadCPREventParameter{}
	for _, param := range requiredCPRParams {
		paramMap[param] = SimPadCPREventParameter{}
	}
	return paramMap
}
func createLogParamMap() map[string]string {
	paramMap := map[string]string{}
	for _, param := range requiredLogParams {
		paramMap[param] = ""
	}
	return paramMap
}

func extractRequiredCPRParams(data *SimPadData) map[string]SimPadCPREventParameter {
	paramMap := createCPRParamMap()

	for _, event := range data.EventList.Events {
		for _, param := range event.Params {
			if _, ok := paramMap[param.Type]; ok {
				if param.Type == "endTime" && event.Type != "CprSessionInfo" {
					continue
				}
				paramMap[param.Type] = param
			}
		}
	}
	return paramMap
}

func extractRequiredLogParams(data *SimPadData) map[string]string {

	var instructor string
	var group string

	if len(data.Log.Instructors.Persons) < 1 {
		instructor = "unknown"
	} else {
		instructor = data.Log.Instructors.Persons[0].Name
	}

	if len(data.Log.Students.Persons) < 3 {
		group = "unknown"
		fmt.Println("Unknown Group here is some known data: ")
		fmt.Printf("Szenario: %v \n", data.Log.Description)
		fmt.Printf("Prüfer: %v \n", instructor)
		fmt.Printf("Timestamp: %v \n", data.Log.SessionDateTimeUTC)
		fmt.Println("Keep in mind that the timestamp is in Coordinated Universal Time")
		fmt.Println("Please add a group and press enter: ")

		var scan string
		_, _ = fmt.Scanln(&scan)
		if scan != "" {
			group = scan
			fmt.Printf("Changed Group to: %v", scan)
		}

	} else {
		group = data.Log.Students.Persons[2].Name
	}

	paramMap := map[string]string{
		"Szenario":  data.Log.Description,
		"Prüfer":    instructor,
		"Gruppe":    group,
		"Fall":      "",
		"Timestamp": data.Log.SessionDateTimeUTC,
	}

	return paramMap
}

func setLogValues(file *excelize.File, logValues map[string]string, key string, col int, row int) {
	paramMap := createLogParamMap()
	if _, ok := paramMap[key]; ok {
		cellName, err := excelize.CoordinatesToCellName(col+1, row+2)
		if err != nil {
			log.Println(err.Error())
			return
		}
		value := logValues[key]
		err = file.SetCellDefault("Daten Kontrolle", cellName, value)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func setCPRValues(file *excelize.File, cprValues map[string]SimPadCPREventParameter, key string, col int, row int) {
	paramMap := createCPRParamMap()
	if _, ok := paramMap[key]; ok {
		cellName, err := excelize.CoordinatesToCellName(col+1, row+2)
		if err != nil {
			log.Println(err.Error())
			return
		}
		value := cprValues[key].Value

		err = file.SetCellDefault("Daten Kontrolle", cellName, value)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func setValuesForEntry(file *excelize.File, data *ProcessedSimPadData) error {

	rows, err := file.GetRows("Daten Kontrolle")
	if err != nil {
		return err
	}

	cols, err := file.GetCols("Daten Kontrolle")
	if err != nil {
		return err
	}

	// New Func
	var firstFreeCol int
	for i, col := range cols[3:] {
		if col[1] == "" {
			firstFreeCol = i + 3
			break
		}
	}

	for r, row := range rows[1:] {
		key := row[1]

		setLogValues(file, data.Log, key, firstFreeCol, r)
		setCPRValues(file, data.CPR, key, firstFreeCol, r)
	}
	return nil
}

func sortData(data []SimPadData) []ProcessedSimPadData {
	dataToDateMap := make(map[string][]ProcessedSimPadData)
	for _, entry := range data {
		parsedTime, err := time.Parse("2006-01-02T15:04:05", entry.Log.SessionDateTimeUTC)
		if err != nil {
			fmt.Println("Could not parse time")
			// TODO: Handle better
			continue
		}
		dateOnly := parsedTime.Format("2006-01-02")

		extractedCPR := extractRequiredCPRParams(&entry)
		extractedLog := extractRequiredLogParams(&entry)
		processed := ProcessedSimPadData{
			extractedLog,
			extractedCPR,
		}

		if list, ok := dataToDateMap[dateOnly]; ok {
			list = append(list, processed)
			dataToDateMap[dateOnly] = list
		} else {
			dataToDateMap[dateOnly] = []ProcessedSimPadData{processed}
		}
	}

	var sortedByGroup [][]ProcessedSimPadData
	for _, list := range dataToDateMap {
		sort.Slice(list, func(i, j int) bool {
			return list[i].Log["Gruppe"] < list[j].Log["Gruppe"]
		})
		sortedByGroup = append(sortedByGroup, list)
	}

	sort.Slice(sortedByGroup, func(i, j int) bool {
		if len(sortedByGroup[i]) < 1 || len(sortedByGroup[j]) < 1 {
			return true // TODO: check if this make sense
		}
		return sortedByGroup[i][0].Log["Timestamp"] < sortedByGroup[j][0].Log["Timestamp"]
	})

	var sorted []ProcessedSimPadData
	for _, list := range sortedByGroup {
		sorted = append(sorted, list...)
	}

	return sorted
}

func parse(inputDir string) error {

	results, err := readAndParseFiles(inputDir)
	if err != nil {
		return err
	}

	f, err := excelize.OpenFile("template.xlsx")
	if err != nil {
		return err
	}

	sorted := sortData(results)

	for i := range sorted {
		err = setValuesForEntry(f, &sorted[i])
		if err != nil {
			return err
		}
	}

	err = f.SaveAs("output.xlsx")
	return err
}
