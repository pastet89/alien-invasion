package world

import (
		"fmt"
		"strings"
)

func (m Map) getEncodedMap() string {
	var lines []string
	mainCityDelimiter := "\n"
	CityRoadsDelimiter := " "	
	for cityName, city := range m.Cities {
		if !city.IsMainCityWithRoadsOut {
			continue
		}
		var line string
		cityRoadsData := city.DestinationCitiesDirections
		line += cityName
		for destination, direction := range cityRoadsData {
			line += CityRoadsDelimiter + direction + "=" + destination
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, mainCityDelimiter)
}

func (m Map) PrintMap() {
	fmt.Println(m.getEncodedMap())
}