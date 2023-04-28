package main

import "encoding/xml"

type SimPadData struct {
	EventList SimPadCPREventList
	Log       SimPadLog
}

type SimPadCPREventParameter struct {
	Text  string `xml:",chardata"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
	Units string `xml:"units,attr"`
}

type SimPadCPREvent struct {
	Text   string                    `xml:",chardata"`
	Type   string                    `xml:"type,attr"`
	Msecs  string                    `xml:"msecs,attr"`
	Params []SimPadCPREventParameter `xml:"param"`
}

type SimPadCPREventList struct {
	XMLName xml.Name         `xml:"eventList"`
	Text    string           `xml:",chardata"`
	Xsi     string           `xml:"xsi,attr"`
	Xsd     string           `xml:"xsd,attr"`
	Events  []SimPadCPREvent `xml:"evt"`
}

type SimPadLog struct {
	XMLName            xml.Name `xml:"Log"`
	Text               string   `xml:",chardata"`
	GeneratedBy        string   `xml:"GeneratedBy,attr"`
	SessionDateTimeUTC string   `xml:"SessionDateTime_UTC,attr"`
	Description        string   `xml:"Description"`
	ScenarioName       string   `xml:"ScenarioName"`
	ScenarioLanguage   string   `xml:"ScenarioLanguage"`
	SessionLanguage    string   `xml:"SessionLanguage"`
	Instructors        struct {
		Text    string `xml:",chardata"`
		Persons []struct {
			Text string `xml:",chardata"`
			Name string `xml:"Name,attr"`
		} `xml:"Person"`
	} `xml:"Instructors"`
	Students struct {
		Text    string `xml:",chardata"`
		Persons []struct {
			Text string `xml:",chardata"`
			Name string `xml:"Name,attr"`
		} `xml:"Person"`
	} `xml:"Students"`
}

type ProcessedSimPadData struct {
	Log map[string]string
	CPR map[string]SimPadCPREventParameter
}
