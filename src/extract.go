package main

import "fmt"

func extractRequiredCPRParams(data *SimPadData) map[string]map[string]SimPadCPREventParameter {
	extractedParams := make(map[string]map[string]SimPadCPREventParameter)

	for _, event := range data.EventList.Events {
		for _, param := range event.Params {
			if _, ok := requiredCPRParams[event.Type][param.Type]; ok {
				if paramMap, _ := extractedParams[event.Type]; paramMap == nil {
					extractedParams[event.Type] = make(map[string]SimPadCPREventParameter)
				}
				extractedParams[event.Type][param.Type] = param
			}
		}
	}
	return extractedParams
}

func extractRequiredLogParams(data *SimPadData) map[string]string {
	scenario := data.LogOverride.Scenario
	instructor := data.LogOverride.Instructor
	group := data.LogOverride.Group
	caseDesc := data.LogOverride.Case

	if instructor == "" && len(data.Log.Instructors.Persons) < 1 {
		instructor = "unknown"
	} else if instructor == "" {
		instructor = data.Log.Instructors.Persons[0].Name
	}

	if group == "" && len(data.Log.Students.Persons) < 3 {
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

	} else if group == "" {
		group = data.Log.Students.Persons[2].Name
	}

	if scenario == "" {
		scenario = data.Log.Description
	}

	paramMap := map[string]string{
		"Szenario":  scenario,
		"Prüfer":    instructor,
		"Gruppe":    group,
		"Fall":      caseDesc,
		"Timestamp": data.Log.SessionDateTimeUTC,
	}

	return paramMap
}
