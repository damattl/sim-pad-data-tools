package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path"
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

		var logOverride = LogOverrideData{}
		logOverrideData, err := os.ReadFile(inputDir + "/" + entry.Name() + "/log-override.json")
		if err == nil {
			err := json.Unmarshal(logOverrideData, &logOverride)
			if err != nil {
				log.Println(err.Error())
			}
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
			EventList:   parsedList,
			Log:         parsedLog,
			LogOverride: logOverride,
		})

	}
	return results, err
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

func parse(inputDir string, execDir string) error {

	results, err := readAndParseFiles(inputDir)
	if err != nil {
		return err
	}

	f, err := excelize.OpenFile(path.Join(execDir, "template.xlsx"))
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

	err = f.SaveAs(path.Join(execDir, "output.xlsx"))
	return err
}
