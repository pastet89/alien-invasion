package world

import (
	"os"
    "strings"
    "github.com/pastet89/alien-invasion/utils"
    "github.com/pastet89/alien-invasion/city"
)

func (m *Map) getMapString() string {
	data, err := os.ReadFile(m.MapFilePath)
	utils.ProcessError(err)
	return string(data)
}

func (m *Map) ParseMap() {
	m.Cities = make(city.Cities)
	m.ReverseRoadsToCities = make(city.ReverseRoadsToCities)
	cities := strings.Split(m.getMapString(), "\n")
	m.parseCities(cities)
}

func (m *Map) parseCitiesRoadsData(cityName string, cityRoadsData []string) {
		for _, cityRoadData := range cityRoadsData {
			splittedCityRoadData := strings.Split(cityRoadData, "=")
			directionName, directionCity := splittedCityRoadData[0], splittedCityRoadData[1]
			isMainCityWithRoadsOut := false
			m.addCity(directionCity, isMainCityWithRoadsOut)
			destinations := m.Cities.GetAllDestinations(cityName)
			destinations[directionCity] = directionName
			m.ReverseRoadsToCities[directionCity] = append(m.ReverseRoadsToCities[directionCity], cityName)
	}
}

func (m *Map) parseCities(cities []string) {
	for _, citiesDataStr := range cities {
		citiesData := strings.Split(citiesDataStr, " ")
		cityName, cityRoadsData := citiesData[0], citiesData[1:]
		isMainCityWithRoadsOut := true
		m.addCity(cityName, isMainCityWithRoadsOut)
		m.parseCitiesRoadsData(cityName, cityRoadsData)
	}
}

func (m *Map) addCity(cityName string, isMainCityWithRoadsOut bool) {
		_, ok := m.Cities[cityName]
		if !ok {
			m.addNonExistingOnMapCity(cityName, isMainCityWithRoadsOut)
		} else {
			m.setExistingOnMapCityFlag(cityName, isMainCityWithRoadsOut)
		}
}

func (m *Map) addNonExistingOnMapCity (cityName string, isMainCityWithRoadsOut bool) {
		destinations := make(city.DestinationCitiesDirections)
		city := &city.City{
					DestinationCitiesDirections: destinations,
					IsMainCityWithRoadsOut: isMainCityWithRoadsOut,
		}
		m.Cities[cityName] = city
}

func (m *Map) setExistingOnMapCityFlag(cityName string, isMainCityWithRoadsOut bool) {
		if isMainCityWithRoadsOut {
			m.Cities[cityName].IsMainCityWithRoadsOut = isMainCityWithRoadsOut
		}
}