package main

import "encoding/xml"

type PIRAlarmXML struct {
	XMLName xml.Name `xml:"PIRAlarm"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Enabled string   `xml:"enabled"`
	Name    string   `xml:"name"`
}

type PIRAlarmStatus struct {
	Name    string
	Enabled int
}

func ParsePIRAlarmStatus(input PIRAlarmXML) (PIRAlarmStatus, error) {
	result := PIRAlarmStatus{}
	result.Name = input.Name
	switch input.Enabled {
	case "true":
		result.Enabled = 1
	case "false":
		result.Enabled = 0
	default:
		result.Enabled = -1
	}
	return result, nil
}

func (s *PIRAlarmStatus) String() string {
	var str string
	switch s.Enabled {
	case 0:
		str = s.Name + " enabled"
	case 1:
		str = s.Name + " disabled"
	default:
		str = s.Name + " error"
	}
	return str
}
