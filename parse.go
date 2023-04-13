package main

import (
	"encoding/xml"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
)

const logPath = "/EventLog.xml"
const cprEventsPath = "/CPR/CPREvents.xml"

func readAndParseFiles() ([]SimPadData, error) {
	entries, err := os.ReadDir("./input")
	if err != nil {
		return []SimPadData{}, err
	}

	results := make([]SimPadData, len(entries))

	for index, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		rawCPRData, err := os.ReadFile("./input/" + entry.Name() + cprEventsPath)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		rawLogData, err := os.ReadFile("./input/" + entry.Name() + logPath)
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
		results[index] = SimPadData{
			EventList: parsedList,
			Log:       parsedLog,
		}
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
	paramMap := map[string]string{
		"Szenario":  data.Log.Description,
		"Prüfer":    data.Log.Instructors.Persons[0].Name,
		"Gruppe":    data.Log.Students.Persons[2].Name,
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

func setValuesForEntry(file *excelize.File, data *SimPadData) {
	extractedCPR := extractRequiredCPRParams(data)
	extractedLog := extractRequiredLogParams(data)

	rows, err := file.GetRows("Daten Kontrolle")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(rows[0][0])

	cols, err := file.GetCols("Daten Kontrolle")
	if err != nil {
		log.Fatal(err)
	}

	// New Func
	var firstFreeCol int
	for i, col := range cols[3:] {
		if col[1] == "" {
			firstFreeCol = i + 3
			break
		}
	}
	log.Println(firstFreeCol)

	for r, row := range rows[1:] {
		key := row[1]

		setLogValues(file, extractedLog, key, firstFreeCol, r)
		setCPRValues(file, extractedCPR, key, firstFreeCol, r)
	}

}

func parse() {
	results, err := readAndParseFiles()
	if err != nil {
		log.Fatal()
	}

	f, err := excelize.OpenFile("template.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	for i := range results {
		setValuesForEntry(f, &results[i])
	}

	err = f.SaveAs("output.xlsx")
	log.Println(err)
}
